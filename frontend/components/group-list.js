class GroupList extends HTMLElement {
  connectedCallback() {
    this.render();
    this.load();
  }

  get phone() {
    return localStorage.getItem('phone');
  }

  async load() {
    const res = await fetch(`/groups?phone=${encodeURIComponent(this.phone)}`);
    if (res.ok) {
      const groups = await res.json();
      this.showGroups(groups);
    } else {
      this.innerHTML = '<p>Failed to load groups</p>';
    }
  }

  showGroups(groups) {
    const list = groups.map(g => `<li data-id="${g.id}">${g.name}</li>`).join('');
    this.querySelector('#group-list').innerHTML = list || '<li>No groups</li>';
  }

  render() {
    this.innerHTML = `
      <h2>Your Groups</h2>
      <ul id="group-list"></ul>
      <button id="create-group">Create Group</button>
    `;
    this.querySelector('#create-group').addEventListener('click', () => {
      this.dispatchEvent(new CustomEvent('create-group', { bubbles: true }));
    });
  }
}

customElements.define('group-list', GroupList);
