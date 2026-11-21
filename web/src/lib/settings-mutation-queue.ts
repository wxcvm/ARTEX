/**
 * Serialize updates to /settings: each response is a full settings snapshot.
 * Pending client writes are overlaid on snapshots so an older response cannot
 * visually undo a newer toggle that has not reached the server yet.
 */
export function createSettingsMutationQueue<T extends object>(
  persist: (patch: Partial<T>) => Promise<T>,
) {
  let tail: Promise<unknown> = Promise.resolve();
  let nextId = 0;
  let pending: Array<{ id: number; patch: Partial<T> }> = [];

  return {
    save(patch: Partial<T>): Promise<T> {
      const id = ++nextId;
      pending.push({ id, patch });
      const write = tail.then(() => persist(patch));
      // A rejected request must not poison the queue for later changes.
      tail = write.then(() => undefined, () => undefined);
      return write.finally(() => {
        pending = pending.filter((entry) => entry.id !== id);
      });
    },
    refresh(read: () => Promise<T>): Promise<T> {
      const request = tail.then(read);
      // Keep later writes behind this read, so recovery cannot observe a
      // snapshot from before an already queued mutation.
      tail = request.then(() => undefined, () => undefined);
      return request.then((snapshot) => this.overlay(snapshot));
    },
    overlay(snapshot: T): T {
      return Object.assign({}, snapshot, ...pending.map(({ patch }) => patch));
    },
  };
}
