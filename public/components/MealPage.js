export class MealPage extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.render();
  }

  render() {
    const homePageTemplate = document.getElementById("meal-page-template");
    console.log(homePageTemplate);
    const templateContent = homePageTemplate.content.cloneNode(true);
    this.appendChild(templateContent);
  }
}

customElements.define("meal-page", MealPage);
