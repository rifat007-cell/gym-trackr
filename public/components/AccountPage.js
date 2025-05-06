import { setTitle } from "../app.js";

export default class AccountPage extends HTMLElement {
  connectedCallback() {
    setTitle("Account");
    const template = document.getElementById("account-page-template");
    const content = template.content.cloneNode(true);
    this.appendChild(content);
    this.render();
  }

  render() {
    const email = app.store.email;

    const h2 = this.querySelector("#email");

    h2.innerHTML = `Welcome <span style = "color:#16B07D;">${email}</span>`;
  }
}
customElements.define("account-page", AccountPage);
