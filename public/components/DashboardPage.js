import { API } from "../services/API.js";

export class DashboardPage extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.render();
  }

  async render() {
    const homePageTemplate = document.getElementById("dashboard-page-template");
    const templateContent = homePageTemplate.content.cloneNode(true);
    this.appendChild(templateContent);

    const data = (await API.getDashboardData()).volumes;

    const labels = data.map((entry) =>
      new Date(entry.log_date).toLocaleDateString()
    );
    const volumes = data.map((entry) => entry.total_volume);

    const ctx = this.querySelector("#dashboard").getContext("2d");

    // Get CSS variables from the page so i use the root variables
    const styles = getComputedStyle(document.documentElement);
    const primaryColor = styles.getPropertyValue("--clr-primary-400").trim();
    const bgColor = styles.getPropertyValue("--clr-bg-main").trim();
    const gridColor = styles.getPropertyValue("--clr-gray-800").trim();
    const fontColor = styles.getPropertyValue("--clr-white").trim();

    new Chart(ctx, {
      type: "line",
      data: {
        labels: labels,
        datasets: [
          {
            label: "Workout Volume",
            data: volumes,
            borderColor: primaryColor,
            backgroundColor: `${primaryColor}33`, // semi-transparent fill
            tension: 0.3,
            fill: true,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            labels: {
              color: fontColor,
              font: {
                family: styles.getPropertyValue("--ff-body").trim(),
              },
            },
          },
        },
        scales: {
          x: {
            ticks: {
              color: fontColor,
              font: {
                family: styles.getPropertyValue("--ff-body").trim(),
              },
            },
            grid: {
              color: gridColor,
            },
            title: {
              display: true,
              text: "Date",
              color: fontColor,
              font: {
                family: styles.getPropertyValue("--ff-heading").trim(),
                size: 14,
              },
            },
          },
          y: {
            ticks: {
              color: fontColor,
              font: {
                family: styles.getPropertyValue("--ff-body").trim(),
              },
            },
            grid: {
              color: gridColor,
            },
            title: {
              display: true,
              text: "Volume (sets × reps × weight)",
              color: fontColor,
              font: {
                family: styles.getPropertyValue("--ff-heading").trim(),
                size: 14,
              },
            },
          },
        },
      },
    });
  }
}

customElements.define("dashboard-page", DashboardPage);
