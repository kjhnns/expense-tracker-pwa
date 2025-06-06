import './components/login-view.js';
import './components/group-list.js';
import './components/group-creator.js';

const appEl = document.getElementById('app');

function showLogin() {
  appEl.innerHTML = '<login-view></login-view>';
}

function showGroups() {
  appEl.innerHTML = '<group-list></group-list>';
}

function showCreator() {
  appEl.innerHTML = '<group-creator></group-creator>';
}

async function handleToken(token) {
  try {
    const res = await fetch(`/verify?token=${encodeURIComponent(token)}`);
    if (res.ok) {
      const data = await res.json();
      localStorage.setItem('phone', data.phone);
      localStorage.removeItem('pendingPhone');
    }
  } catch {}
}

document.addEventListener('DOMContentLoaded', async () => {
  const params = new URLSearchParams(window.location.search);
  const token = params.get('token');
  if (token) {
    await handleToken(token);
    history.replaceState({}, '', '/');
  }

  if (localStorage.getItem('phone')) {
    showGroups();
  } else {
    showLogin();
  }
});

appEl.addEventListener('create-group', showCreator);
