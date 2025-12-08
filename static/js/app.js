const API_URL = '/api/v1'
let currentCart = [];

document.addEventListener('DOMContentLoaded', () => {
    checkAuth();
    loadProducts();
})


function checkAuth() {
    const token = localStorage.getItem('token');
    const email = localStorage.getItem('email');
    const role = localStorage.getItem('role');

    if (token) {
        document.getElementById('guest-nav').classList.add('hidden');
        document.getElementById('user-nav').classList.remove('hidden');
        document.getElementById('user-email-display').innerText = email;

        // Show Admin button only if admin
        if (role === 'admin') {
            document.getElementById('admin-btn').classList.remove('hidden');
        }
    } else {
        document.getElementById('guest-nav').classList.remove('hidden');
        document.getElementById('user-nav').classList.add('hidden');
        document.getElementById('admin-panel').classList.add('hidden');
    }
}

// --- Auth Functions ---
async function login() {
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-pass').value;

    try {
        const res = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({email, password})
        });

        const data = await res.json();
        if (!res.ok) throw new Error(data.error);

        // Decode JWT simply to get role (in production use a library)
        // We'll just rely on what we saved, or you can fetch profile
        // For this demo, let's assume we want to save the token
        localStorage.setItem('token', data.token);
        localStorage.setItem('email', email);

        // Hacky way to check admin without decoding JWT on client:
        // Try to hit an admin endpoint or just save it from response if you modified the backend
        // For now, let's just assume if email contains "admin" (simple test) or save manually
        // Better way: Backend sends role in response body alongside token
        if (email.includes("admin")) localStorage.setItem('role', 'admin');
        else localStorage.setItem('role', 'user');

        showToast("Logged in successfully!", "success");
        checkAuth();
        window.location.reload(); // Refresh to clean state
    } catch (err) {
        showToast(err.message, "error");
    }
}

async function register() {
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-pass').value;
    try {
        const res = await fetch(`${API_URL}/auth/register`, {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify({email, password})
        });
        const data = await res.json();
        if (!res.ok) throw new Error(data.error);
        showToast("Account created! Please login.", "success");
    } catch (err) {
        showToast(err.message, "error");
    }
}

function logout() {
    localStorage.clear();
    window.location.reload();
}

function toggleAdminPanel() {
    const panel = document.getElementById('admin-panel');
    panel.classList.toggle('hidden');
}

// --- Product Logic ---
async function loadProducts() {
    try {
        const res = await fetch(`${API_URL}/products`);
        const products = await res.json();
        const container = document.getElementById('product-list');
        container.innerHTML = '';

        if (products.length === 0) {
            container.innerHTML = '<div class="col-span-full text-center text-gray-500">No games found. Admin needs to add stock!</div>';
            return;
        }

        products.forEach(p => {
            const price = (p.price / 100).toFixed(2);
            // Random gradients for game covers since we don't have real images
            const colors = ['from-purple-500 to-indigo-500', 'from-green-500 to-teal-500', 'from-red-500 to-orange-500'];
            const bg = colors[p.ID % colors.length];

            const html = `
                    <div class="bg-gray-800 rounded-xl overflow-hidden shadow-lg card-hover border border-gray-700 flex flex-col">
                        <div class="h-40 bg-gradient-to-r ${bg} flex items-center justify-center">
                            <span class="text-4xl">🎮</span>
                        </div>
                        <div class="p-5 flex-grow">
                            <div class="flex justify-between items-start">
                                <h3 class="text-xl font-bold text-white mb-2">${p.name}</h3>
                                <span class="bg-gray-700 text-xs px-2 py-1 rounded text-gray-300">Stock: ${p.stock}</span>
                            </div>
                            <p class="text-gray-400 text-sm mb-4 line-clamp-2">${p.description || 'No description available.'}</p>
                            <div class="flex justify-between items-center mt-auto">
                                <span class="text-2xl font-bold text-green-400">$${price}</span>
                                <button onclick="buyProduct(${p.ID})" 
                                    class="${p.stock > 0 ? 'bg-blue-600 hover:bg-blue-500' : 'bg-gray-600 cursor-not-allowed'} text-white px-4 py-2 rounded-lg font-bold transition">
                                    ${p.stock > 0 ? 'Buy Now' : 'Sold Out'}
                                </button>
                            </div>
                        </div>
                    </div>
                    `;
            container.insertAdjacentHTML('beforeend', html);
        });
    } catch (err) {
        console.error(err);
        showToast("Failed to load products", "error");
    }
}

async function addProduct() {
    const token = localStorage.getItem('token');
    if (!token) return showToast("You must be logged in!", "error");

    const product = {
        name: document.getElementById('p-name').value,
        description: document.getElementById('p-desc').value,
        price: parseInt(document.getElementById('p-price').value),
        stock: parseInt(document.getElementById('p-stock').value),
        sku: document.getElementById('p-sku').value
    };

    try {
        const res = await fetch(`${API_URL}/products`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify(product)
        });

        const data = await res.json();
        if (!res.ok) throw new Error(data.error);

        showToast("Product added successfully!", "success");
        loadProducts(); // Refresh list
    } catch (err) {
        showToast(err.message, "error");
    }
}

async function buyProduct(id) {
    const token = localStorage.getItem('token');
    if (!token) {
        showToast("Please login to purchase games", "error");
        return;
    }

    try {
        const res = await fetch(`${API_URL}/orders`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({product_id: id})
        });

        const data = await res.json();
        if (!res.ok) throw new Error(data.error);

        showToast("🎉 " + data.message, "success");
        loadProducts(); // Refresh to show reduced stock
    } catch (err) {
        showToast(err.message, "error");
    }
}

// --- Helper: Toast Notification ---
function showToast(message, type = "success") {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    const bg = type === "success" ? "bg-green-600" : "bg-red-600";

    toast.className = `${bg} text-white px-6 py-3 rounded-lg shadow-xl flex items-center gap-3 transform transition-all duration-300 translate-y-10 opacity-0`;
    toast.innerHTML = `<span>${type === "success" ? "✅" : "⚠️"}</span> <span>${message}</span>`;

    container.appendChild(toast);

    // Animate in
    setTimeout(() => toast.classList.remove('translate-y-10', 'opacity-0'), 10);

    // Remove after 3 seconds
    setTimeout(() => {
        toast.classList.add('translate-y-10', 'opacity-0');
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}
