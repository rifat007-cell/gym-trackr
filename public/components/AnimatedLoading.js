export class AnimatedLoading extends HTMLElement {
  constructor() {
    super();
  }

  // when we change the data-elements, data-width, or data-height attributes in the inspector manually.
  // we want to re-render the component
  static get observedAttributes() {
    return ["data-elements", "data-width", "data-height"];
  }

  connectedCallback() {
    this.render();
  }

  attributeChangedCallback(name, oldValue, newValue) {
    if (oldValue !== newValue) {
      this.render();
    }
  }

  render() {
    this.innerHTML = ""; // Clear previous elements

    const elements = Number(this.dataset.elements) || 5; // Default to 5
    const width = this.dataset.width || "20px";
    const height = this.dataset.height || "20px";

    for (let i = 0; i < elements; i++) {
      const wrapper = document.createElement("div");
      wrapper.classList.add("loading-wave");
      wrapper.style.width = width;
      wrapper.style.height = height;
      wrapper.style.margin = "10px";
      wrapper.style.display = "inline-block";
      this.appendChild(wrapper);
    }
  }
}

customElements.define("animated-loading", AnimatedLoading);
