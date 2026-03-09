const outputEl = document.getElementById("output");

function getBaseUrl() {
  return document.getElementById("baseUrl").value.trim().replace(/\/$/, "");
}

function showResponse(label, status, data) {
  outputEl.textContent = `${label}\nstatus: ${status}\n\n${JSON.stringify(data, null, 2)}`;
}

async function apiRequest(path, options = {}) {
  const base = getBaseUrl();
  const url = `${base}${path}`;
  const resp = await fetch(url, {
    headers: { "Content-Type": "application/json" },
    ...options
  });

  let body;
  try {
    body = await resp.json();
  } catch {
    body = { raw: await resp.text() };
  }
  return { status: resp.status, body };
}

document.getElementById("loadProductsBtn").addEventListener("click", async () => {
  try {
    const { status, body } = await apiRequest("/api/products");
    showResponse("GET /api/products", status, body);
  } catch (err) {
    showResponse("GET /api/products", "network error", { error: String(err) });
  }
});

document.getElementById("loginForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const login = document.getElementById("login").value.trim();
  const password = document.getElementById("password").value;

  try {
    const { status, body } = await apiRequest("/api/login", {
      method: "POST",
      body: JSON.stringify({ login, password })
    });
    showResponse("POST /api/login", status, body);
  } catch (err) {
    showResponse("POST /api/login", "network error", { error: String(err) });
  }
});

document.getElementById("orderForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const contactData = document.getElementById("contactData").value.trim();
  const comment = document.getElementById("comment").value.trim();
  const productID = Number(document.getElementById("productId").value);
  const quantity = Number(document.getElementById("quantity").value);

  try {
    const { status, body } = await apiRequest("/api/orders", {
      method: "POST",
      body: JSON.stringify({
        contact_data: contactData,
        comment,
        products: [{ product_id: productID, quantity }]
      })
    });
    showResponse("POST /api/orders", status, body);
  } catch (err) {
    showResponse("POST /api/orders", "network error", { error: String(err) });
  }
});
