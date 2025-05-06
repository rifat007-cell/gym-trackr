export class HomePage extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.render();
  }

  async render() {
    const homePageTemplate = document.getElementById("home-page-template");
    console.log(homePageTemplate);
    const templateContent = homePageTemplate.content.cloneNode(true);
    this.appendChild(templateContent);

    this.querySelector("a").addEventListener("click", (e) => {
      e.preventDefault();
      const href = e.target.getAttribute("href");
      app.router.go(href);
    });
  }
}

customElements.define("home-page", HomePage);
