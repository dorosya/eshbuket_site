const $ = (sel) => document.querySelector(sel);

async function api(path, options = {}) {
  const res = await fetch(path, options);
  const contentType = res.headers.get("content-type") || "";
  let data = null;
  if (contentType.includes("application/json")) {
    data = await res.json().catch(() => null);
  } else {
    data = await res.text().catch(() => null);
  }
  if (!res.ok) {
    const msg = (data && data.error) ? data.error : (data || `HTTP ${res.status}`);
    throw new Error(msg);
  }
  return data;
}

function setMsg(el, text, isErr = false) {
  el.textContent = text || "";
  el.classList.toggle("err", !!isErr);
}

function productCard(p) {
  const div = document.createElement("div");
  div.className = "product";
  const img = document.createElement("img");
  img.alt = p.name;
  img.loading = "lazy";
  img.src = `/api/products/${p.id}/image`;
  img.onerror = () => {
    img.src = "data:image/svg+xml;charset=utf-8," + encodeURIComponent(
      `<svg xmlns='http://www.w3.org/2000/svg' width='400' height='260'>
        <rect width='100%' height='100%' fill='#f2f2f2'/>
        <text x='50%' y='50%' dominant-baseline='middle' text-anchor='middle' fill='#777' font-size='18'>
          no image
        </text>
      </svg>`
    );
  };

  div.innerHTML = `
    <div class="product__meta">
      <div class="product__name"></div>
      <div class="muted small">ID: ${p.id}</div>
      <div class="muted small">Категория: ${p.category || "-"}</div>
      <div class="product__price">${p.price}</div>
    </div>
  `;
  div.prepend(img);
  div.querySelector(".product__name").textContent = p.name;
  return div;
}

async function loadProducts() {
  const listMsg = $("#listMsg");
  setMsg(listMsg, "");
  const wrap = $("#products");
  wrap.innerHTML = "";
  const select = $("#productSelect");
  select.innerHTML = "";

  try {
    const items = await api("/api/products");
    if (!Array.isArray(items) || items.length === 0) {
      setMsg(listMsg, "Пока нет товаров. Создай один выше 🙂");
      return;
    }

    for (const p of items) {
      wrap.appendChild(productCard(p));

      const opt = document.createElement("option");
      opt.value = p.id;
      opt.textContent = `${p.id} — ${p.name}`;
      select.appendChild(opt);
    }
  } catch (e) {
    setMsg(listMsg, `Ошибка загрузки: ${e.message}`, true);
  }
}

async function onCreateProduct(e) {
  e.preventDefault();
  const msg = $("#createProductMsg");
  setMsg(msg, "");

  const form = e.currentTarget;
  const payload = {
    name: form.name.value.trim(),
    price: form.price.value.trim(),
    category: form.category.value.trim(),
  };

  try {
    await api("/api/products", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload),
      // NOTE: если включишь auth — тут нужно будет выставлять cookie сессии
    });
    setMsg(msg, "Товар создан ✅");
    form.reset();
    await loadProducts();
  } catch (e2) {
    setMsg(msg, `Ошибка: ${e2.message}`, true);
  }
}

async function onUploadImage(e) {
  e.preventDefault();
  const msg = $("#uploadMsg");
  setMsg(msg, "");

  const form = e.currentTarget;
  const productId = form.productId.value;
  const file = form.image.files[0];
  if (!productId || !file) return;

  const fd = new FormData();
  fd.append("image", file);

  try {
    await api(`/api/products/${productId}/image`, {
      method: "POST",
      body: fd,
    });
    setMsg(msg, "Загружено ✅");
    form.reset();
    await loadProducts();
  } catch (e2) {
    setMsg(msg, `Ошибка: ${e2.message}`, true);
  }
}

document.addEventListener("DOMContentLoaded", () => {
  $("#createProductForm").addEventListener("submit", onCreateProduct);
  $("#uploadImageForm").addEventListener("submit", onUploadImage);
  $("#refreshBtn").addEventListener("click", loadProducts);
  loadProducts();
});
