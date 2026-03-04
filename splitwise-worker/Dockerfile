FROM python:3.13-slim

ENV UV_COMPILE_BYTECODE=1 \
    UV_LINK_MODE=copy \
    UV_TOOL_BIN_DIR=/usr/local/bin \
    \
    PIP_NO_CACHE_DIR=off \
    PIP_DISABLE_PIP_VERSION_CHECK=on \
    PIP_DEFAULT_TIMEOUT=100 \
    \
    PYSETUP_PATH="/opt/pysetup" \
    VENV_PATH="/opt/pysetup/.venv"

ENV PATH="$VENV_PATH/bin:$PATH"

WORKDIR $PYSETUP_PATH

COPY --from=ghcr.io/astral-sh/uv:0.9.1 /uv /uvx /bin/

COPY pyproject.toml uv.lock ./

RUN uv venv $VENV_PATH && \
    uv sync --locked --no-install-project

WORKDIR /app

COPY . .

CMD ["python", "-m", "src.main"]