package data

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"time"

	"github.com/tanvir-rifat007/gymBuddy/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
    ErrDuplicateEmail = errors.New("duplicate email")
)



type User struct{
	ID int `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Name string `json:"name"`
	Email string `json:"email"`
	Password Password `json:"-"`
	Activated bool `json:"activated"`
	Version int `json:"-"`
  JWT string `json:"jwt,omitempty"`
	
}

type Password struct{
	Plaintext *string 
	hash []byte

}

type UserModel struct{
	DB *sql.DB
}

func (p *Password) Set(plaintext string) error{
	hash,err:=bcrypt.GenerateFromPassword([]byte(plaintext),bcrypt.DefaultCost)

	if err!=nil{
		return err
	}
	p.Plaintext = &plaintext
	p.hash = hash

	

	return nil
}

func (p *Password) Matches(plaintext string) (bool,error){
    err:= bcrypt.CompareHashAndPassword(p.hash,[]byte(plaintext))

	if err!=nil{
		if errors.Is(err,bcrypt.ErrMismatchedHashAndPassword){
			return false,nil
		}
		return false,err

	}
	return true,nil
}

func (m UserModel) Insert(user *User) error{
	stmt:= `INSERT INTO users (name,email,password_hash,activated)
	        VALUES($1,$2,$3,$4) RETURNING id,created_at,version`

	args:= []any{user.Name,user.Email,user.Password.hash,user.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	err:= m.DB.QueryRowContext(ctx,stmt,args...).Scan(&user.ID,&user.CreatedAt,&user.Version)

	  if err != nil {
        switch {
        case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
            return ErrDuplicateEmail
        default:
            return err
        }
			}

	return nil
}


func (m UserModel) GetUserByEmail(email string) (*User,error){
	stmt:= `SELECT id,created_at,name,email,password_hash,activated,version
          FROM users WHERE email = $1`

	var user User
	ctx,cancel:=context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()
	err:= m.DB.QueryRowContext(ctx,stmt,email).Scan(&user.ID,&user.CreatedAt,&user.Name,&user.Email,&user.Password.hash,&user.Activated,&user.Version)

	if err!=nil{
		if err==sql.ErrNoRows{
			return nil,ErrRecordNotFound
		}else{
			return nil,err
		}
	}

	return &user,nil


}

func (m UserModel) Update(user *User) error {
    query := `
        UPDATE users 
        SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
        WHERE id = $5 AND version = $6
        RETURNING version`

    args := []any{
        user.Name,
        user.Email,
        user.Password.hash,
        user.Activated,
        user.ID,
        user.Version,
    }

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    err := m.DB.QueryRowContext(ctx, query, args...).Scan(&user.Version)
    if err != nil {
        switch {
        case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
            return ErrDuplicateEmail
        case errors.Is(err, sql.ErrNoRows):
            return ErrEditConflict
        default:
            return err
        }
    }

    return nil
}

func (m UserModel) Login(email, password string) (*User, error) {
    // Get the user with the specified email address.
    user, err := m.GetUserByEmail(email)
    if err != nil {
        return nil, err
    }

    // Check if the provided password matches the stored hash.
    matched, err := user.Password.Matches(password)
    if err != nil {
        return nil, err
    }
    if !matched {
        return nil, ErrInvalidCredentials
    }

    // Return the matching user.
    return user, nil
}



func (m UserModel) GetUserFromToken(tokenScope, tokenPlaintext string) (*User, error) {
    // Calculate the SHA-256 hash of the plaintext token provided by the client.
    // Remember that this returns a byte *array* with length 32, not a slice.
    tokenHash := sha256.Sum256([]byte(tokenPlaintext))

    // Set up the SQL query.
    query := `
        SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = $1
        AND tokens.scope = $2 
        AND tokens.expiry > $3`

    // Create a slice containing the query arguments. Notice how we use the [:] operator
    // to get a slice containing the token hash, rather than passing in the array (which
    // is not supported by the pq driver), and that we pass the current time as the
    // value to check against the token expiry.
    args := []any{tokenHash[:], tokenScope, time.Now()}

    var user User

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Execute the query, scanning the return values into a User struct. If no matching
    // record is found we return an ErrRecordNotFound error.
    err := m.DB.QueryRowContext(ctx, query, args...).Scan(
        &user.ID,
        &user.CreatedAt,
        &user.Name,
        &user.Email,
        &user.Password.hash,
        &user.Activated,
        &user.Version,
    )
    if err != nil {
        switch {
        case errors.Is(err, sql.ErrNoRows):
            return nil, ErrRecordNotFound
        default:
            return nil, err
        }
    }

    // Return the matching user.
    return &user, nil
}


func (m *UserModel) LogReminderSent(userID int) error {
    stmt := `INSERT INTO reminder_logs (user_id) VALUES ($1) ON CONFLICT DO NOTHING`

    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    _, err := m.DB.ExecContext(ctx, stmt, userID)
    return err
}


func (m *UserModel) GetUsersMissingWorkoutLogs() ([]User, error) {
    query := `
        SELECT u.id, u.email
        FROM users u
        LEFT JOIN (
            SELECT user_id, MAX(log_date) AS last_log_date
            FROM user_workout_logs
            GROUP BY user_id
        ) l ON u.id = l.user_id
        WHERE
            (l.last_log_date IS NULL OR l.last_log_date < CURRENT_DATE - INTERVAL '1 day')
            AND u.activated = true
            AND NOT EXISTS (
                SELECT 1 FROM reminder_logs r
                WHERE r.user_id = u.id AND r.reminder_date = CURRENT_DATE
            );
    `

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    rows, err := m.DB.QueryContext(ctx, query)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var u User
        if err := rows.Scan(&u.ID, &u.Email); err != nil {
            return nil, err
        }
        users = append(users, u)
    }

    return users, rows.Err()
}




func ValidateEmail(v *validator.Validator, email string) {
    v.Check(email != "", "email", "must be provided")
    v.Check(validator.Matches(email, validator.EmailRX), "email", "must be a valid email address")
}

func ValidatePasswordPlaintext(v *validator.Validator, password string) {
    v.Check(password != "", "password", "must be provided")
    v.Check(len(password) >= 8, "password", "must be at least 8 bytes long")
    v.Check(len(password) <= 72, "password", "must not be more than 72 bytes long")
}


func ValidateUser(v *validator.Validator, user *User) {
    v.Check(user.Name != "", "name", "must be provided")
    v.Check(len(user.Name) <= 500, "name", "must not be more than 500 bytes long")

    // Call the standalone ValidateEmail() helper.
    ValidateEmail(v, user.Email)

    // If the plaintext password is not nil, call the standalone 
    // ValidatePasswordPlaintext() helper.
    if user.Password.Plaintext != nil {
        ValidatePasswordPlaintext(v, *user.Password.Plaintext)
    }

    // If the password hash is ever nil, this will be due to a logic error in our 
    // codebase (probably because we forgot to set a password for the user). It's a 
    // useful sanity check to include here, but it's not a problem with the data 
    // provided by the client. So rather than adding an error to the validation map we 
    // raise a panic instead.
    if user.Password.hash == nil {
        panic("missing password hash for user")
    }
}