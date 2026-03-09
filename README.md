[![Tests](https://github.com/dorosya/eshbuket_site/actions/workflows/tests.yaml/badge.svg)](https://github.com/dorosya/eshbuket_site/actions/workflows/tests.yaml)

## Local frontend smoke run

1. Start backend on `http://localhost:8080`.
2. Run a static server for frontend:
   - `cd frontend/static`
   - `python -m http.server 5500`
3. Open `http://localhost:5500` and use the buttons/forms to call:
   - `POST /api/login`
   - `GET /api/products`
   - `POST /api/orders`
