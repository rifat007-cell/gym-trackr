import AccountPage from "../components/AccountPage.js";
import { ActivatedPage } from "../components/ActivatedPage.js";
import { DashboardPage } from "../components/DashboardPage.js";
import { HomePage } from "../components/HomePage.js";
import { LoginPage } from "../components/LoginPage.js";
import { MealPage } from "../components/MealPage.js";
import { RegisterPage } from "../components/RegisterPage.js";
import { WorkoutLogPage } from "../components/WorkoutLogPage.js";
import { WorkoutPage } from "../components/WorkoutPage.js";

export const routes = [
  {
    path: "/",
    component: HomePage,
  },
  {
    path: "/workout",
    component: WorkoutPage,
    loggedIn: true,
  },

  {
    path: "/meal",
    component: MealPage,
    loggedIn: true,
  },

  {
    path: "/account/register",
    component: RegisterPage,
  },

  {
    path: "/account/login",
    component: LoginPage,
  },

  {
    path: "/account/",
    component: AccountPage,
    loggedIn: true,
  },

  {
    path: "/activated",
    component: ActivatedPage,
  },

  {
    path: "/workoutlog",
    component: WorkoutLogPage,
    loggedIn: true,
  },

  {
    path: "/dashboard",
    component: DashboardPage,
    loggedIn: true,
  },
];
