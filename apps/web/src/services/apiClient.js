const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";

export class ApiRequestError extends Error {
  constructor(message, options = {}) {
    super(message);
    this.name = "ApiRequestError";
    this.status = options.status;
    this.code = options.code || "request_failed";
    this.hint = options.hint || "";
    this.details = options.details || "";
    this.job = options.job || null;
    this.payload = options.payload || null;
  }
}

async function parseResponse(response) {
  const contentType = response.headers.get("content-type") || "";
  const isJSON = contentType.includes("application/json");
  const payload = isJSON ? await response.json() : await response.text();

  if (!response.ok) {
    const errorPayload = typeof payload === "string" ? { message: payload } : payload?.error || {};
    const message =
      typeof errorPayload === "string"
        ? errorPayload
        : errorPayload.message || `Request failed with status ${response.status}`;

    throw new ApiRequestError(message.trim(), {
      status: response.status,
      code: typeof errorPayload === "object" ? errorPayload.code : undefined,
      hint: typeof errorPayload === "object" ? errorPayload.hint : undefined,
      details: typeof errorPayload === "object" ? errorPayload.details : undefined,
      job: typeof payload === "object" ? payload?.job : null,
      payload,
    });
  }

  return payload;
}

export async function apiRequest(path, options = {}) {
  const response = await fetch(`${API_BASE_URL}${path}`, options);
  return parseResponse(response);
}