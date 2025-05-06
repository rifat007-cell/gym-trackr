import { setTitle } from "../app.js";

export class WorkoutLogPage extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    setTitle("Workout Log");
    this.render();
  }

  render() {
    const homePageTemplate = document.getElementById(
      "workoutlog-page-template"
    );
    console.log(homePageTemplate);
    const templateContent = homePageTemplate.content.cloneNode(true);
    this.appendChild(templateContent);

    this.style.zIndex = "-1000";
  }
}

customElements.define("workout-log-page", WorkoutLogPage);
