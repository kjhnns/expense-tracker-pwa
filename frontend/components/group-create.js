class GroupCreate extends HTMLElement {
  connectedCallback() {
    this.innerHTML = `<button>Create Group</button>`;
  }
}
customElements.define('group-create', GroupCreate);
