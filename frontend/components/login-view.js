class LoginView extends HTMLElement {
  connectedCallback() {
    this.render();
  }

  render() {
    this.innerHTML = `
      <form id="login-form">
        <label>Phone Number
          <input type="tel" id="phone" placeholder="+123456789" required />
        </label>
        <button type="submit">Send Login Link</button>
      </form>
      <p id="message"></p>
    `;
    this.querySelector("#login-form").addEventListener("submit", async (e) => {
      e.preventDefault();
      const phone = this.querySelector("#phone").value.trim();
      try {
        const res = await fetch("/register", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ phone }),
        });
        if (!res.ok) throw new Error("Request failed");
        const data = await res.json();
        localStorage.setItem("pendingPhone", phone);
        this.querySelector("#message").innerHTML =
          `Open <a href="${data.link}">this link</a> to verify your login.`;
      } catch (err) {
        this.querySelector("#message").textContent = "Failed to start login";
      }
    });
  }
}

customElements.define("login-view", LoginView);
