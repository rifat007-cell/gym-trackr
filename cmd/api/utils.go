package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"os"
	"strings"
	"text/template"

	"github.com/tanvir-rifat007/gymBuddy/internal/data"
)


type envelope map[string]any


func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {

    js, err := json.MarshalIndent(data, "", "\t")
    if err != nil {
        return err
    }

    js = append(js, '\n')


    for key, value := range headers {
        w.Header()[key] = value
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    w.Write(js)

    return nil
}


func (app *application) readJSON(w http.ResponseWriter, r *http.Request,dst any)error{
	// limit the req body to 1mb:
	maxBytes:= 1_048_576;
	r.Body = http.MaxBytesReader(w,r.Body,int64(maxBytes))
	dec:= json.NewDecoder(r.Body)

	// this will make sure that the json decoder returns an error if the request body contains any additional fields which cannot be mapped to the target destination
	dec.DisallowUnknownFields()

	err:= dec.Decode(dst)


	if err!=nil{
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError


		switch{
		case errors.As(err,&syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)",syntaxError.Offset)

		case errors.Is(err,io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
            if unmarshalTypeError.Field != "" {
                return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
            }
            return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)


    case errors.Is(err, io.EOF):
            return errors.New("body must not be empty")

	   case strings.HasPrefix(err.Error(), "json: unknown field "):
            fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
            return fmt.Errorf("body contains unknown key %s", fieldName)

        // Use the errors.As() function to check whether the error has the type 
        // *http.MaxBytesError. If it does, then it means the request body exceeded our 
        // size limit of 1MB and we return a clear error message.
        case errors.As(err, &maxBytesError):
            return fmt.Errorf("body must not be larger than %d bytes", maxBytesError.Limit)



    case errors.As(err, &invalidUnmarshalError):
            panic(err)

    default:
            return err
		}


	}

	// if there is any additional JSON data in the request body, return an error
// like : curl -d '{"title": "Moana"}{"title": "Top Gun"}' localhost:4000/v1/movies

// here for this 2nd one {"title": "Top Gun"} we will get an error
	 err = dec.Decode(&struct{}{})
    if !errors.Is(err, io.EOF) {
        return errors.New("body must only contain a single JSON value")
    }
return nil
}


// i am using the built in go smtp package to send emails

// func (app *application) sendEmail(to []string, subject string, data any) error {
// 	auth := smtp.PlainAuth("", os.Getenv("FROM_EMAIL"), os.Getenv("FROM_EMAIL_PASSWORD"), os.Getenv("FROM_EMAIL_SMTP"))

// 	// 1. Parse the template
// 	tmpl, err := template.ParseFiles("./internal/mailer/templates/user_welcome.tmpl.html")
// 	if err != nil {
// 		return err
// 	}

// 	// 2. Render the template to a buffer
// 	buf:= new(bytes.Buffer)
// 	err = tmpl.Execute(buf, data)
// 	if err != nil {
// 		return err
// 	}

// 	// 3. Construct email headers
// 	headers := make(map[string]string)
// 	headers["From"] = os.Getenv("FROM_EMAIL")
// 	headers["To"] = strings.Join(to, ",")
// 	headers["Subject"] = subject
// 	headers["MIME-Version"] = "1.0"
// 	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

// 	var message strings.Builder
// 	for k, v := range headers {
// 		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
// 	}
// 	message.WriteString("\r\n" + buf.String())

// 	// 4. Send the email
// 	return smtp.SendMail(
// 		os.Getenv("SMTP_ADDR"), // e.g., "sandbox.smtp.mailtrap.io:587"
// 		auth,
// 		os.Getenv("FROM_EMAIL"),
// 		to,
// 		[]byte(message.String()),
// 	)
// }


func (app *application) sendEmail(to []string, subject string, templateFile string, data any) error {
	auth := smtp.PlainAuth("", os.Getenv("FROM_EMAIL"), os.Getenv("FROM_EMAIL_PASSWORD"), os.Getenv("FROM_EMAIL_SMTP"))

	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, data)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"From":         os.Getenv("FROM_EMAIL"),
		"To":           strings.Join(to, ","),
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": "text/html; charset=\"UTF-8\"",
	}

	var message strings.Builder
	for k, v := range headers {
		message.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	message.WriteString("\r\n" + buf.String())

	return smtp.SendMail(
		os.Getenv("SMTP_ADDR"),
		auth,
		os.Getenv("FROM_EMAIL"),
		to,
		[]byte(message.String()),
	)
}




func (app *application )SendWorkoutReminderEmails(users data.UserModel) {
	userList, err := users.GetUsersMissingWorkoutLogs()
	if err != nil {
		app.logger.Error("Failed to fetch users for reminder", "error", err)
		return
	}

	for _, u := range userList {
		subject := "💪 Time to log your workout!"
		// body := "Hey! It looks like you missed logging your workout yesterday. Keep the streak going — log it now."

		var data struct{
			Name string
		}

		data.Name = u.Name



		err := app.sendEmail([]string{u.Email}, subject, "./internal/mailer/templates/workout_reminder.tmpl.html", data)
		if err != nil {
			app.logger.Error("Failed to send reminder", "user", u.ID, "error", err)
			continue
		}

		err = users.LogReminderSent(u.ID)
		if err != nil {
			app.logger.Error("Failed to log reminder sent", "user", u.ID, "error", err)
		} else {
			app.logger.Info("Reminder sent", "user", u.ID)
		}
	}
}
