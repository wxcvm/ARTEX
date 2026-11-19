/** A credential check can fail even while the current session is still valid. */
export function shouldClearSession(
  path: string,
  status: number,
  message: string,
  requestToken: string | null,
  currentToken: string | null,
): boolean {
  if (status !== 401 || requestToken !== currentToken) return false;
  if (path === "/auth/login") return false;
  return path !== "/auth/change-password" || message !== "Current password is wrong";
}

/** Ignore late model-status replies belonging to a session that went away. */
export async function resolveSignInDestination(
  token: string,
  currentToken: () => string | null,
  loadModel: () => Promise<{ configured: boolean }>,
): Promise<"/function/tasks" | "/welcome" | null> {
  let configured = false;
  try {
    configured = (await loadModel()).configured;
  } catch {
    // Other model-status failures should not prevent sign-in.
  }
  if (currentToken() !== token) return null;
  return configured ? "/function/tasks" : "/welcome";
}
