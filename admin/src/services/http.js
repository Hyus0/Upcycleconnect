export class ApiError extends Error {
  constructor(message, status = 500, payload = null) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.payload = payload;
  }
}

export class MissingEndpointError extends Error {
  constructor(resource, detail) {
    super(detail ?? `Endpoint manquant pour ${resource}`);
    this.name = "MissingEndpointError";
    this.resource = resource;
  }
}

async function parseResponseBody(response) {
  const contentType = response.headers.get("content-type") ?? "";
  if (contentType.includes("application/json")) {
    return response.json();
  }
  return response.text();
}

export async function request(url, options = {}) {
  const response = await fetch(url, {
    headers: {
      "Content-Type": "application/json",
      ...(options.headers ?? {})
    },
    ...options
  });

  const payload = await parseResponseBody(response);

  if (!response.ok) {
    const message =
      typeof payload === "string"
        ? payload
        : payload?.error ?? payload?.message ?? `Erreur HTTP ${response.status}`;
    throw new ApiError(message, response.status, payload);
  }

  return payload;
}
