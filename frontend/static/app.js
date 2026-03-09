const outputEl = document.getElementById("output");
const imagePreviewEl = document.getElementById("imagePreview");

function getBaseUrl() {
  const raw = document.getElementById("baseUrl").value.trim().replace(/\/$/, "");
  if (raw) {
    return raw;
  }

  const host = window.location.hostname;
  const port = window.location.port;
  if ((host === "localhost" || host === "127.0.0.1") && port === "5500") {
    return "http://localhost:8080";
  }

  return "";
}

function getFullUrl(path) {
  return `${getBaseUrl()}${path}`;
}

function showResponse(label, status, data) {
  outputEl.textContent = `${label}\nstatus: ${status}\n\n${JSON.stringify(data, null, 2)}`;
}

function showError(label, err) {
  showResponse(label, "network error", { error: String(err) });
}

async function apiRequest(path, options = {}) {
  const url = getFullUrl(path);
  const headers = options.body instanceof FormData ? {} : { "Content-Type": "application/json" };
  const resp = await fetch(url, {
    credentials: "include",
    headers,
    ...options
  });

  const contentType = resp.headers.get("content-type") || "";
  let body;
  if (contentType.includes("application/json")) {
    body = await resp.json();
  } else {
    body = { raw: await resp.text() };
  }

  return { status: resp.status, body, ok: resp.ok };
}

async function runSmokeTest() {
  try {
    const products = await apiRequest("/api/products");
    showResponse("SMOKE: GET /api/products", products.status, products.body);
  } catch (err) {
    showError("SMOKE: GET /api/products", err);
  }
}

document.getElementById("smokeBtn").addEventListener("click", runSmokeTest);

document.getElementById("productsFilterForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const category = document.getElementById("productsCategory").value.trim();
  const path = category ? `/api/products?category=${encodeURIComponent(category)}` : "/api/products";

  try {
    const { status, body } = await apiRequest(path);
    showResponse("GET /api/products", status, body);
  } catch (err) {
    showError("GET /api/products", err);
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
    showError("POST /api/login", err);
  }
});

document.getElementById("createProductForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const name = document.getElementById("productName").value.trim();
  const price = document.getElementById("productPrice").value.trim();
  const category = document.getElementById("productCategory").value.trim();

  try {
    const { status, body } = await apiRequest("/api/products", {
      method: "POST",
      body: JSON.stringify({ name, price, category })
    });
    showResponse("POST /api/products", status, body);
  } catch (err) {
    showError("POST /api/products", err);
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
    showError("POST /api/orders", err);
  }
});

document.getElementById("uploadImageForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const productID = document.getElementById("uploadImageProductId").value.trim();
  const fileInput = document.getElementById("productImageFile");
  const file = fileInput.files && fileInput.files[0];

  if (!file) {
    showResponse("POST /api/products/:id/image", 400, { error: "image file is required" });
    return;
  }

  const form = new FormData();
  form.append("image", file);

  try {
    const { status, body } = await apiRequest(`/api/products/${productID}/image`, {
      method: "POST",
      body: form
    });
    showResponse("POST /api/products/:id/image", status, body);
  } catch (err) {
    showError("POST /api/products/:id/image", err);
  }
});

document.getElementById("viewImageForm").addEventListener("submit", async (e) => {
  e.preventDefault();
  const productID = document.getElementById("viewImageProductId").value.trim();
  const url = getFullUrl(`/api/products/${productID}/image`);

  imagePreviewEl.src = url;
  imagePreviewEl.classList.remove("hidden");
  imagePreviewEl.onload = () => showResponse("GET /api/products/:id/image", 200, { image_url: url });
  imagePreviewEl.onerror = async () => {
    imagePreviewEl.classList.add("hidden");
    try {
      const { status, body } = await apiRequest(`/api/products/${productID}/image`);
      showResponse("GET /api/products/:id/image", status, body);
    } catch (err) {
      showError("GET /api/products/:id/image", err);
    }
  };
});
