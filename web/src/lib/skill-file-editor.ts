// The selection object is the identity of one editor visit, even for A -> B -> A.
export type FileSelection = { skill: string; path: string };
export type FileEditorState = {
  selection: FileSelection;
  status: "loading" | "ready" | "error";
  content: string;
  dirty: boolean;
  saving: boolean;
};

export function beginFileRead(selection: FileSelection): FileEditorState {
  return { selection, status: "loading", content: "", dirty: false, saving: false };
}

export function updateFileEditor(
  previous: FileEditorState | null,
  selection: FileSelection,
  patch: Partial<Omit<FileEditorState, "selection">>,
): FileEditorState | null {
  return previous?.selection === selection ? { ...previous, ...patch } : previous;
}
