export default class AccountPage extends HTMLElement {
  connectedCallback() {
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
