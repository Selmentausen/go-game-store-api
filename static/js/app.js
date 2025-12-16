const API_URL = '/api/v1';
let currentCart = [];

// --- Init on Page Load ---
document.addEventListener('DOMContentLoaded', () => {
    checkAuth();
    loadProducts();
});

// --- Auth Functions ---
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

        // Load the cart count
        fetchCart();
    } else {
        document.getElementById('guest-nav').classList.remove('hidden');
        document.getElementById('user-nav').classList.add('hidden');
        document.getElementById('admin-panel').classList.add('hidden');
    }
}

async function login() {
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-pass').value;

    try {
        const res = await fetch(`${API_URL}/auth/login`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
        });

        const data = await res.json();
        if (!res.ok) throw new Error(data.error);

        localStorage.setItem('token', data.token);
        localStorage.setItem('email', email);

        // Simple role check (in prod, parse the JWT)
        if(email.includes("admin")) localStorage.setItem('role', 'admin');
        else localStorage.setItem('role', 'user');

        showToast("Logged in successfully!", "success");
        checkAuth();
        window.location.reload();
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
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
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

// --- Product Functions ---
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
            // Random gradients for visuals
            const colors = ['from-purple-500 to-indigo-500', 'from-green-500 to-teal-500', 'from-red-500 to-orange-500'];
            const bg = colors[p.ID % colors.length];

            const html = `
            <div class="bg-gray-800 rounded-xl overflow-hidden shadow-lg card-hover border border-gray-700 flex flex-col">
                <div class="h-40 bg-gradient-to-r ${bg} flex items-center justify-center">
                    <span class="text-4xl">🎮</span>
                </div>
                <div class="p-5 flex-grow flex flex-col">
                    <div class="flex justify-between items-start">
                        <h3 class="text-xl font-bold text-white mb-2">${p.name}</h3>
                        <span class="bg-gray-700 text-xs px-2 py-1 rounded text-gray-300">Stock: ${p.stock}</span>
                    </div>
                    <p class="text-gray-400 text-sm mb-4 line-clamp-2">${p.description || 'No description available.'}</p>
                    <div class="flex justify-between items-center mt-auto">
                        <span class="text-2xl font-bold text-green-400">$${price}</span>
                        
                        <!-- UPDATED BUTTON: Calls addToCart instead of buyProduct -->
                        <button onclick="addToCart(${p.ID})" 
                            class="${p.stock > 0 ? 'bg-blue-600 hover:bg-blue-500' : 'bg-gray-600 cursor-not-allowed'} text-white px-4 py-2 rounded-lg font-bold transition">
                            ${p.stock > 0 ? 'Add to Cart' : 'Sold Out'}
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

// --- Cart Logic (The New Part) ---

async function fetchCart() {
    const token = localStorage.getItem('token');
    if (!token) return;

    try {
        const res = await fetch(`${API_URL}/cart`, {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        if (res.ok) {
            currentCart = await res.json();
            updateCartUI();
        }
    } catch (e) { console.error(e); }
}

async function addToCart(id) {
    const token = localStorage.getItem('token');
    if (!token) return showToast("Please login to shop", "error");

    try {
        const res = await fetch(`${API_URL}/cart`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ product_id: id, quantity: 1 })
        });

        if (!res.ok) throw new Error((await res.json()).error);

        showToast("Added to Cart", "success");
        fetchCart(); // Update Badge and UI
    } catch (err) {
        showToast(err.message, "error");
    }
}

function updateCartUI() {
    const badge = document.getElementById('cart-badge');
    const list = document.getElementById('cart-items');
    const totalEl = document.getElementById('cart-total');
    const btn = document.getElementById('checkout-btn');

    // Update Badge Count
    const totalItems = currentCart.reduce((sum, item) => sum + item.quantity, 0);
    badge.innerText = totalItems;
    badge.classList.toggle('hidden', totalItems === 0);

    // Calculate Total Price
    let totalCents = 0;

    // Render List
    if (currentCart.length === 0) {
        list.innerHTML = '<div class="text-center text-gray-500 mt-10">Your cart is empty.</div>';
        btn.disabled = true;
        totalEl.innerText = "$0.00";
        return;
    }

    btn.disabled = false;
    list.innerHTML = '';

    currentCart.forEach(item => {
        // Ensure product data exists
        if(item.product) {
            totalCents += item.product.price * item.quantity;
            list.innerHTML += `
            <div class="flex justify-between items-center bg-gray-700 p-3 rounded-lg border border-gray-600 transition hover:bg-gray-600">
                <div class="flex-grow">
                    <div class="font-bold text-sm text-white">${item.product.name}</div>
                    <div class="text-xs text-gray-400">$${(item.product.price/100).toFixed(2)} each</div>
                </div>
                
                <div class="flex items-center gap-3">
                    <!-- QUANTITY CONTROLS -->
                    <div class="flex items-center bg-gray-800 rounded">
                        <button onclick="window.changeQuantity(${item.product_id}, -1)" class="px-2 py-1 text-gray-300 hover:text-white hover:bg-gray-600 rounded-l">-</button>
                        <span class="px-2 text-sm font-mono">${item.quantity}</span>
                        <button onclick="window.changeQuantity(${item.product_id}, 1)" class="px-2 py-1 text-gray-300 hover:text-white hover:bg-gray-600 rounded-r">+</button>
                    </div>

                    <span class="font-mono font-bold text-green-400 w-16 text-right">$${(item.product.price * item.quantity / 100).toFixed(2)}</span>
                    
                    <!-- Full Remove (Trash) -->
                    <button onclick="window.removeFromCart(${item.product_id})" class="text-red-400 hover:text-red-200 p-1">✕</button>
                </div>
            </div>`;
        }
    });

    totalEl.innerText = "$" + (totalCents/100).toFixed(2);
}

async function checkout() {
    const token = localStorage.getItem('token');
    if (!token) return;

    if(!confirm("Confirm purchase?")) return;

    try {
        const res = await fetch(`${API_URL}/cart/checkout`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            }
        });

        const data = await res.json();
        if (!res.ok) throw new Error(data.error);

        showToast(`🎉 Success! Order #${data.order_id} placed.`, "success");

        // Refresh everything
        toggleCart(); // Close modal
        fetchCart();  // Clear cart UI
        loadProducts(); // Update stock on main page
    } catch (err) {
        showToast(err.message, "error");
    }
}

// --- UI Helpers ---
function toggleCart() {
    const modal = document.getElementById('cart-modal');
    const panel = document.getElementById('cart-panel');

    if (modal.classList.contains('hidden')) {
        modal.classList.remove('hidden');
        setTimeout(() => panel.classList.remove('translate-x-full'), 10);
    } else {
        panel.classList.add('translate-x-full');
        setTimeout(() => modal.classList.add('hidden'), 300);
    }
}

function toggleAdminPanel() {
    document.getElementById('admin-panel').classList.toggle('hidden');
}

async function addProduct() {
    const token = localStorage.getItem('token');
    if (!token) return;

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

        if (!res.ok) throw new Error((await res.json()).error);

        showToast("Product added successfully!", "success");
        loadProducts();
    } catch (err) {
        showToast(err.message, "error");
    }
}

async function removeFromCart(productID) {
    const token = localStorage.getItem('token');
    if (!token) return;

    try {
        const res = await fetch(`${API_URL}/cart/${productID}`, {
            method: 'DELETE',
            headers: {
                "Authorization": `Bearer ${token}`
            }
        });
        if (!res.ok) throw new Error("Failed to remove item");

        fetchCart();
        showToast("Item removed", "success");
    } catch (err) {
        showToast(err.message, "error");
    }
}

async function changeQuantity(productID, delta) {
    const token = localStorage.getItem('token');
    if (!token) return;

    try {
        const res = await fetch(`${API_URL}/cart/`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({ product_id: productID, quantity: delta }),
        });
        if (!res.ok) throw new Error((await res.json()).error);
        fetchCart();
    } catch (err) {
        showToast(err.message, "error");
    }
}

function showToast(message, type = "success") {
    const container = document.getElementById('toast-container');
    const toast = document.createElement('div');
    const bg = type === "success" ? "bg-green-600" : "bg-red-600";

    toast.className = `${bg} text-white px-6 py-3 rounded-lg shadow-xl flex items-center gap-3 transform transition-all duration-300 translate-y-10 opacity-0`;
    toast.innerHTML = `<span>${type === "success" ? "✅" : "⚠️"}</span> <span>${message}</span>`;

    container.appendChild(toast);

    setTimeout(() => toast.classList.remove('translate-y-10', 'opacity-0'), 10);

    setTimeout(() => {
        toast.classList.add('translate-y-10', 'opacity-0');
        setTimeout(() => toast.remove(), 300);
    }, 3000);
}

// Global scope attachments for HTML onClick events
window.login = login;
window.register = register;
window.logout = logout;
window.addToCart = addToCart;
window.checkout = checkout;
window.toggleCart = toggleCart;
window.toggleAdminPanel = toggleAdminPanel;
window.addProduct = addProduct;
window.removeFromCart = removeFromCart;
window.changeQuantity = changeQuantity;