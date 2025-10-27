import { backendBaseUrl } from "@/services/baseUrl";

export function resolveImageUrl(path, fallback = null) {
  if (!path) return fallback;
  if (/^https?:\/\//.test(path)) return path;
  if (path.startsWith("/")) return backendBaseUrl + path;
  return backendBaseUrl + "/" + path;
}
