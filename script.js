document.addEventListener('DOMContentLoaded', (event) => {
    fetch('/.netlify/functions/api/links')
        .then(res => res.json())
        .then(links => {
            const linksList = document.getElementById('links-list');
            document.getElementById('loading-indicator').style.display = 'none';
            links.forEach(link => {
                linksList.innerHTML += `
                <div class="bg-white p-4 rounded-lg shadow-sm">
                    <a href="${link.url}" target="_blank" class="text-blue-600 font-medium hover:underline">${link.title || link.url}</a>
                    <p class="mt-2 text-gray-700">${link.commentary}</p>
                </div>`;
            });
        });
});

document.body.addEventListener('htmx:beforeSwap', (event) => {
    if (event.detail.target.id === 'links-list') {
        document.getElementById('form-loading-indicator').classList.remove('hidden');
    }
});

document.body.addEventListener('htmx:afterSwap', (event) => {
    if (event.detail.target.id === 'links-list') {
        document.getElementById('form-loading-indicator').classList.add('hidden');
        event.detail.target.scrollIntoView({ behavior: 'smooth' });
        event.detail.target.closest('form').reset();
    }
    if (event.detail.target.id === 'auth-section') {
        if (event.detail.xhr.status === 200) {
            document.getElementById('auth-section').classList.add('hidden');
            document.getElementById('add-link-section').innerHTML = ""; // Clear existing content
            htmx.ajax('GET', '/add-link', '#add-link-section');
            document.getElementById('add-link-section').classList.remove('hidden');
        } else {
            const errorMessage = JSON.parse(event.detail.xhr.response).message
            alert(errorMessage)
        }
    }
});

document.addEventListener('DOMContentLoaded', () => {
    htmx.ajax('GET', '/login', '#login-form');
});
