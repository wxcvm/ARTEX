export const REPORT_PAGE_SIZE = 50;

export type ReportLibraryLoad<T> = {
  key: string;
  revision: number;
  status: "loading" | "ready" | "error";
  items: T[];
  error?: string;
};

export function reportLibraryView<T>(
  load: ReportLibraryLoad<T> | null,
  key: string,
  revision: number,
) {
  const sameFilterPage = load?.key === key;
  const current = sameFilterPage && load.revision === revision;
  const status = current ? load.status : "loading";
  const items = sameFilterPage ? load.items : [];
  return {
    items,
    status,
    error: current && load.status === "error" ? load.error : undefined,
    stale: items.length > 0 && status !== "ready",
    empty: status === "ready" && items.length === 0,
    canAct: status === "ready",
  };
}

export function reportOffsetAfterDeletion(offset: number, rows: number) {
  return rows === 1 ? Math.max(0, offset - REPORT_PAGE_SIZE) : offset;
}
