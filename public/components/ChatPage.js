import { API } from "../services/API.js";

export class ChatWidget extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: "open" });

    shadow.innerHTML = `
<style>
  :host { all: initial; }

  #chat-toggle-btn {
    position: fixed;
    bottom: 20px;
    right: 20px;
    background: var(--clr-green-500);
    color: var(--clr-white);
    border: none;
    border-radius: 50%;
    width: 56px;
    height: 56px;
    font-size: 1.5rem;
    cursor: pointer;
    z-index: 10000;
    display: flex;
    align-items: center;
    justify-content: center;
    box-shadow: 0 4px 15px rgba(0, 0, 0, 0.2);
    transition: background 0.2s ease;
  }

  #chat-toggle-btn:hover {
    background: var(--clr-green-600);
  }

  #chat-widget {
    position: fixed;
    bottom: 90px;
    right: 20px;
    width: 320px;
    max-height: 500px;
    background: var(--clr-brown-700);
    border-radius: 1rem;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
    overflow: hidden;
    display: flex;
    flex-direction: column;
    font-family: var(--ff-body);
    z-index: 9999;
    transition: transform 0.3s ease, opacity 0.3s ease;
  }

  #chat-widget.hidden {
    transform: translateY(20px);
    opacity: 0;
    pointer-events: none;
  }

  #chat-header {
    background: var(--clr-green-500);
    color: var(--clr-white);
    padding: 0.75rem 1rem;
    font-weight: bold;
    font-size: var(--fs-500);
    text-align: center;
    font-family: var(--ff-heading);
  }

  #chat-body {
    flex: 1;
    padding: 0.625rem;
    overflow-y: auto;
    background: var(--clr-brown-600);
    font-size: var(--fs-400);
    display: flex;
    flex-direction: column;
  }

  .message {
    margin: 0.5rem 0;
    padding: 0.625rem 0.875rem;
    border-radius: 0.75rem;
    line-height: 1.4;
    max-width: 80%;
    font-size: var(--fs-400);
  }

  .message.bot {
    background: var(--clr-gray-100);
    color: var(--clr-brown-900);
    align-self: flex-start;
  }

  .message.user {
    background: var(--clr-green-400);
    color: var(--clr-white);
    align-self: flex-end;
  }

  #chat-form {
    display: flex;
    border-top: 1px solid var(--clr-gray-100);
    background: var(--clr-brown-700);
  }

  #chat-input {
    flex: 1;
    border: none;
    padding: 0.75rem;
    outline: none;
    font-size: var(--fs-400);
    font-family: var(--ff-body);
    color: var(--clr-white);
    background: var(--clr-brown-800);
    border-radius: 0 0 0 1rem;
  }

  #chat-form button {
    background: var(--clr-green-500);
    color: var(--clr-white);
    border: none;
    padding: 0 1rem;
    cursor: pointer;
    font-size: var(--fs-500);
    font-family: var(--ff-body);
    border-radius: 0 0 1rem 0;
    transition: background 0.2s ease-in-out;
  }

  #chat-form button:hover {
    background: var(--clr-green-600);
  }
</style>

      <button id="chat-toggle-btn" title="Toggle Chat">💬</button>

      <div id="chat-widget" class="hidden">
        <div id="chat-header">💪 AI Coach</div>
        <div id="chat-body">
          <div class="message bot">Hey! I'm your fitness coach. Ask me anything!</div>
        </div>
        <form id="chat-form">
          <input type="text" id="chat-input" placeholder="Type your message..." />
          <button type="submit">➤</button>
        </form>
      </div>
    `;

    const form = shadow.querySelector("#chat-form");
    const input = shadow.querySelector("#chat-input");
    const body = shadow.querySelector("#chat-body");
    const widget = shadow.querySelector("#chat-widget");
    const toggleBtn = shadow.querySelector("#chat-toggle-btn");

    this.historyLoaded = false;

    form.addEventListener("submit", (e) => {
      e.preventDefault();
      const message = input.value.trim();
      if (!message) return;

      const userMsg = document.createElement("div");
      userMsg.className = "message user";
      userMsg.textContent = message;
      body.appendChild(userMsg);
      input.value = "";

      const loadingMsg = document.createElement("div");
      loadingMsg.className = "message bot";
      loadingMsg.textContent = "Let me think...";
      body.appendChild(loadingMsg);
      body.scrollTop = body.scrollHeight;

      const data = { message };

      async function fetchData() {
        try {
          const res = await API.chatWithAI(data);
          body.removeChild(loadingMsg);

          const botMsg = document.createElement("div");
          botMsg.className = "message bot";
          botMsg.textContent =
            res.response ??
            "You have to login and activate your account before chatting with me...";
          body.appendChild(botMsg);
          body.scrollTop = body.scrollHeight;
        } catch (error) {
          console.error("Error:", error);
          loadingMsg.textContent = "An error occurred. Please try again.";
        }
      }

      fetchData();
    });

    toggleBtn.addEventListener("click", () => {
      widget.classList.toggle("hidden");
      toggleBtn.textContent = widget.classList.contains("hidden") ? "💬" : "➖";

      if (!widget.classList.contains("hidden") && !this.historyLoaded) {
        this.loadChatHistory(body);
        this.historyLoaded = true;
      }
    });
  }

  async loadChatHistory(container) {
    try {
      const res = await API.getChatHistory();
      console.log("Respnse", res);
      if (!res.history || res.history.length === 0) return;

      res.history.forEach((msg) => {
        const msgDiv = document.createElement("div");
        msgDiv.className = `message ${msg.Role === "user" ? "user" : "bot"}`;
        msgDiv.textContent = msg.Content;
        container.appendChild(msgDiv);
      });

      container.scrollTop = container.scrollHeight;
    } catch (err) {
      console.error("Failed to load chat history:", err);
    }
  }
}

customElements.define("chat-widget", ChatWidget);
