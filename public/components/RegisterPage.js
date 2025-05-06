import { setTitle } from "../app.js";

export class RegisterPage extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    setTitle("Register");

    this.render();
  }

  render() {
    const registerPageTemplate = document.getElementById(
      "register-page-template"
    );
    console.log(registerPageTemplate);
    const templateContent = registerPageTemplate.content.cloneNode(true);
    this.appendChild(templateContent);

    this.querySelector("a").addEventListener("click", (event) => {
      event.preventDefault();
      app.router.go("/account/login");
    });
  }
}

customElements.define("register-page", RegisterPage);
