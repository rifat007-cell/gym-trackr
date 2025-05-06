import { setTitle } from "../app.js";

export class LoginPage extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    setTitle("Login");

    this.render();
  }

  render() {
    const workoutPageTemplate = document.getElementById("login-page-template");
    console.log(workoutPageTemplate);
    const templateContent = workoutPageTemplate.content.cloneNode(true);
    this.appendChild(templateContent);

    this.querySelector("a").addEventListener("click", (event) => {
      event.preventDefault();
      app.router.go("/account/register");
    });
  }
}

customElements.define("login-page", LoginPage);
