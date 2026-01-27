const API = "http://localhost:8080/api";

async function loadProducts() {
  const res = await fetch(`${API}/products`);
  const products = await res.json();
  const ul = document.getElementById("products");
  ul.innerHTML = "";
  products.forEach(p => {
    const li = document.createElement("li");
    li.innerHTML = `#${p.id} ${p.name} — ${p.price} ₽
      <br><img src="/api/products/${p.id}/image" width="100" onerror="this.style.display='none'"/>`;
    ul.appendChild(li);
  });
}

async function createProduct() {
  const body = {
    name: document.getElementById("name").value,
    price: Number(document.getElementById("price").value),
    category: document.getElementById("category").value
  };
  await fetch(`${API}/products`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body)
  });
  loadProducts();
}

async function uploadImage() {
  const id = document.getElementById("productId").value;
  const file = document.getElementById("image").files[0];
  const form = new FormData();
  form.append("image", file);

  await fetch(`${API}/products/${id}/image`, {
    method: "POST",
    body: form
  });

  document.getElementById("preview").innerHTML =
    `<img src="/api/products/${id}/image" width="200"/>`;
  loadProducts();
}

loadProducts();
