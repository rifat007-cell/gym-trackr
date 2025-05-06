export class WorkoutPage extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.render();
  }

  render() {
    const workoutPageTemplate = document.getElementById(
      "workout-page-template"
    );
    console.log(workoutPageTemplate);
    const templateContent = workoutPageTemplate.content.cloneNode(true);
    this.appendChild(templateContent);
  }
}

customElements.define("workout-page", WorkoutPage);
