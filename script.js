document.addEventListener('DOMContentLoaded', (event) => {
    fetch('/.netlify/functions/api/')
        .then(res => res.json())
        .then(links => {
            const linksList = document.getElementById('links-list');
            document.getElementById('loading-indicator').style.display = 'none';
            links.forEach(link => {
                linksList.innerHTML += `
                <div class="bg-white p-4 rounded-lg shadow-sm">
                    <a href="<span class="math-inline">\{link\.url\}" target\="\_blank" class\="text\-blue\-600 font\-medium hover\:underline"\></span>{link.title || link.url}</a>
                    <p class="mt-2 text-gray-700">${link.commentary}</p>
                </div>`;
            });
        });
});

document.body.addEventListener('htmx:beforeSwap', (
