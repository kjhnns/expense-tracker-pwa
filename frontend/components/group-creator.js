class GroupCreator extends HTMLElement {
  connectedCallback() {
    this.render();
  }

  render() {
    this.innerHTML = `
      <form id="group-form" class="container">
        <label>Group Name
          <input type="text" id="group-name" required />
        </label>
        <label>Default Currency
          <select id="default-currency">
            <option value="EUR">EUR</option>
            <option value="USD">USD</option>
          </select>
        </label>
        <label>Your Phone Number
          <input type="tel" id="creator-phone" placeholder="+123456789" required />
        </label>
        <div id="participants">
          <label>Participant
            <input type="tel" class="participant-phone" placeholder="+123456789" />
          </label>
        </div>
        <button type="button" id="add-participant">Add Participant</button>
        <button type="submit">Create Group</button>
      </form>
      <article id="confirmation" style="display:none"></article>
    `;

    this.querySelector('#add-participant').addEventListener('click', () => {
      const wrapper = document.createElement('label');
      wrapper.textContent = 'Participant';
      const input = document.createElement('input');
      input.type = 'tel';
      input.className = 'participant-phone';
      input.placeholder = '+123456789';
      wrapper.appendChild(document.createElement('br'));
      wrapper.appendChild(input);
      this.querySelector('#participants').appendChild(wrapper);
    });

    this.querySelector('#group-form').addEventListener('submit', async (e) => {
      e.preventDefault();
      const groupName = this.querySelector('#group-name').value.trim();
      const currency = this.querySelector('#default-currency').value;
      const createdBy = this.querySelector('#creator-phone').value.trim();
      const participantInputs = Array.from(this.querySelectorAll('.participant-phone'));
      const others = participantInputs.map(i => i.value.trim()).filter(v => v);
      const participants = [createdBy, ...others];
      const payload = {
        group_name: groupName,
        default_currency: currency,
        created_by: createdBy,
        participants
      };
      try {
        const res = await fetch('/groups/create', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(payload)
        });
        if (!res.ok) throw new Error('Request failed');
        const data = await res.json();
        this.showConfirmation(data);
      } catch (err) {
        this.showError(err);
      }
    });
  }

  showConfirmation(data) {
    const article = this.querySelector('#confirmation');
    const links = Object.entries(data.invite_links)
      .map(([phone, link]) => `<li>${phone}: <a href="${link}">${link}</a></li>`)
      .join('');
    article.innerHTML = `<h3>Group Created</h3>
      <p>ID: ${data.group_id}</p>
      <ul>${links}</ul>`;
    article.style.display = 'block';
  }

  showError(err) {
    const article = this.querySelector('#confirmation');
    article.innerHTML = `<p>Failed to create group: ${err}</p>`;
    article.style.display = 'block';
  }
}

customElements.define('group-creator', GroupCreator);
