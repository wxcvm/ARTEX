"use client";

import { useUiLocale } from "@/lib/ui-locale.client";
import { uiText } from "@/lib/ui-locale";

import * as React from "react";

import { LibraryIcon, PlusIcon, SaveIcon, Settings2Icon, Trash2Icon } from "lucide-react";
import { toast } from "sonner";

import { AssetInterceptRulesEditor } from "@/components/asset-intercept-rules-editor";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import {
  Combobox,
  ComboboxContent,
  ComboboxEmpty,
  ComboboxInput,
  ComboboxItem,
  ComboboxList,
} from "@/components/ui/combobox";
import { Field, FieldDescription, FieldGroup, FieldLabel } from "@/components/ui/field";
import { Input } from "@/components/ui/input";
import { NativeSelect, NativeSelectOption } from "@/components/ui/native-select";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Sheet, SheetContent, SheetDescription, SheetFooter, SheetHeader, SheetTitle } from "@/components/ui/sheet";
import { Spinner } from "@/components/ui/spinner";
import { Textarea } from "@/components/ui/textarea";
import { api } from "@/lib/api";
import type { AssetInterceptRuleInput, TaskCategory, TaskTemplate } from "@/lib/types";
import { cn } from "@/lib/utils";

interface TemplateDraft {
  name: string;
  description: string;
  goal: string;
  categoryID: number | null;
  interceptRules: AssetInterceptRuleInput[];
}

// TemplateSeed holds the initial values carried over from the creation form when using "Save as template".
type TemplateSeed = Pick<TemplateDraft, "description" | "goal" | "categoryID" | "interceptRules">;

interface TaskTemplateManagerProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  templates: TaskTemplate[];
  seed: TemplateSeed | null;
  onCreated: (template: TaskTemplate) => void;
  onUpdated: (template: TaskTemplate) => void;
  onDeleted: (id: number) => void;
}

const emptyDraft = (): TemplateDraft => ({
  name: "",
  description: "",
  goal: "",
  categoryID: null,
  interceptRules: [],
});

function templateDraft(template: TaskTemplate): TemplateDraft {
  return {
    name: template.name,
    description: template.description,
    goal: template.goal,
    categoryID: template.category_id ?? null,
    interceptRules: template.intercept_rules ?? [],
  };
}

function TaskTemplateManager({
  open,
  onOpenChange,
  templates,
  seed,
  onCreated,
  onUpdated,
  onDeleted,
}: TaskTemplateManagerProps) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [selectedID, setSelectedID] = React.useState<number | null>(null);
  const [draft, setDraft] = React.useState<TemplateDraft>(emptyDraft);
  const [saving, setSaving] = React.useState(false);
  const [deleteOpen, setDeleteOpen] = React.useState(false);
  const [deleting, setDeleting] = React.useState(false);
  const [categories, setCategories] = React.useState<TaskCategory[]>([]);
  const wasOpen = React.useRef(false);

  React.useEffect(() => {
    if (!open) return;
    api
      .taskCategories()
      .then(setCategories)
      .catch(() => setCategories([]));
  }, [open]);

  React.useEffect(() => {
    if (open && !wasOpen.current) {
      if (seed) {
        setSelectedID(null);
        setDraft({
          name: "",
          description: seed.description,
          goal: seed.goal,
          categoryID: seed.categoryID,
          interceptRules: seed.interceptRules,
        });
      } else if (templates[0]) {
        setSelectedID(templates[0].id);
        setDraft(templateDraft(templates[0]));
      } else {
        setSelectedID(null);
        setDraft(emptyDraft());
      }
    }
    wasOpen.current = open;
  }, [open, seed, templates]);

  const selectTemplate = (template: TaskTemplate) => {
    setSelectedID(template.id);
    setDraft(templateDraft(template));
  };

  const startNew = () => {
    setSelectedID(null);
    setDraft(emptyDraft());
  };

  const updateDraft = (field: "name" | "description" | "goal", value: string) => {
    setDraft((current) => ({ ...current, [field]: value }));
  };

  const patchDraft = (patch: Partial<TemplateDraft>) => {
    setDraft((current) => ({ ...current, ...patch }));
  };

  async function save() {
    const input = {
      name: draft.name.trim(),
      description: draft.description.trim(),
      goal: draft.goal.trim(),
      category_id: draft.categoryID,
      intercept_rules: draft.interceptRules
        .map((r) => ({ ...r, pattern: r.pattern.trim() }))
        .filter((r) => r.pattern !== ""),
    };
    if (!input.name || !input.description || !input.goal) {
      toast.error(tr("请填写模板名称、描述和目标", "Fill in the template name, description, and goal"));
      return;
    }
    setSaving(true);
    try {
      if (selectedID == null) {
        const created = await api.createTaskTemplate(input);
        onCreated(created);
        setSelectedID(created.id);
        setDraft(templateDraft(created));
        toast.success(tr("模板已创建", "Template created"));
      } else {
        const updated = await api.updateTaskTemplate(selectedID, input);
        onUpdated(updated);
        setDraft(templateDraft(updated));
        toast.success(tr("模板已更新", "Template updated"));
      }
    } catch (error) {
      toast.error(tr(`保存失败：${(error as Error).message}`, `Save failed: ${(error as Error).message}`));
    } finally {
      setSaving(false);
    }
  }

  async function remove() {
    if (selectedID == null) return;
    const deletedID = selectedID;
    setDeleting(true);
    try {
      await api.deleteTaskTemplate(deletedID);
      onDeleted(deletedID);
      const next = templates.find((template) => template.id !== deletedID);
      if (next) {
        setSelectedID(next.id);
        setDraft(templateDraft(next));
      } else {
        startNew();
      }
      setDeleteOpen(false);
      toast.success(tr("模板已删除", "Template deleted"));
    } catch (error) {
      toast.error(tr(`删除失败：${(error as Error).message}`, `Delete failed: ${(error as Error).message}`));
    } finally {
      setDeleting(false);
    }
  }

  let saveLabel = saving ? tr("正在保存", "Saving") : tr("保存修改", "Save changes");
  if (!saving && selectedID == null) saveLabel = tr("创建模板", "Create template");

  return (
    <>
      <Sheet open={open} onOpenChange={onOpenChange}>
        <SheetContent className="grid h-full w-full! max-w-none! grid-rows-[auto_minmax(0,1fr)_auto] gap-0 overflow-hidden p-0 sm:w-[48rem]! sm:max-w-[48rem]!">
          <SheetHeader className="border-b px-6 py-5">
            <SheetTitle>{tr("任务模板管理", "Task template management")}</SheetTitle>
            <SheetDescription>{tr("模板保存任务描述、目标、分类和任务级拦截/放行规则；编辑不会影响已创建的任务。", "Templates store the description, goal, category, and task-level intercept/allow rules; edits do not affect already-created tasks.")}</SheetDescription>
          </SheetHeader>
          <div className="grid min-h-0 overflow-y-auto lg:grid-cols-[15rem_minmax(0,1fr)] lg:overflow-hidden">
            <div className="flex min-h-0 flex-col border-b p-3 lg:border-r lg:border-b-0">
              <Button type="button" variant="outline" className="w-full" onClick={startNew}>
                <PlusIcon data-icon="inline-start" />
                {tr("创建模板", "Create template")}
              </Button>
              <ScrollArea className="mt-2 max-h-44 lg:max-h-none lg:flex-1">
                <div className="flex flex-col gap-1 pr-2">
                  {templates.length === 0 && (
                    <p className="px-2 py-6 text-center text-muted-foreground text-sm">{tr("暂无模板", "No templates yet")}</p>
                  )}
                  {templates.map((template) => (
                    <button
                      key={template.id}
                      type="button"
                      className={cn(
                        "min-w-0 rounded-md px-2.5 py-2 text-left transition-colors",
                        selectedID === template.id ? "bg-accent text-accent-foreground" : "hover:bg-accent/50",
                      )}
                      onClick={() => selectTemplate(template)}
                    >
                      <span className="block truncate font-medium text-sm">{template.name}</span>
                      <span className="block truncate text-muted-foreground text-xs">{template.description}</span>
                    </button>
                  ))}
                </div>
              </ScrollArea>
            </div>
            <ScrollArea className="min-h-0">
              <FieldGroup className="p-6">
                <Field>
                  <FieldLabel htmlFor="task-template-name">{tr("模板名称", "Template name")}</FieldLabel>
                  <Input
                    id="task-template-name"
                    value={draft.name}
                    maxLength={120}
                    placeholder={tr("例如：外部 Web 安全测试", "e.g. External web penetration testing")}
                    onChange={(event) => updateDraft("name", event.target.value)}
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="task-template-description">{tr("描述", "Description")}</FieldLabel>
                  <Textarea
                    id="task-template-description"
                    className="min-h-28"
                    value={draft.description}
                    placeholder={tr("测试目标及背景", "Test target and background")}
                    onChange={(event) => updateDraft("description", event.target.value)}
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="task-template-goal">{tr("目标", "Goal")}</FieldLabel>
                  <Textarea
                    id="task-template-goal"
                    className="min-h-28"
                    value={draft.goal}
                    placeholder={tr("任务需要达到的目标", "Goals the task needs to achieve")}
                    onChange={(event) => updateDraft("goal", event.target.value)}
                  />
                </Field>
                <Field>
                  <FieldLabel htmlFor="task-template-category">{tr("任务分类", "Task category")}</FieldLabel>
                  <NativeSelect
                    id="task-template-category"
                    className="w-full"
                    value={draft.categoryID == null ? "" : String(draft.categoryID)}
                    onChange={(event) =>
                      patchDraft({ categoryID: event.target.value === "" ? null : Number(event.target.value) })
                    }
                  >
                    <NativeSelectOption value="">{tr("未分类", "Uncategorized")}</NativeSelectOption>
                    {categories.map((c) => (
                      <NativeSelectOption key={c.id} value={String(c.id)}>
                        {c.name}
                      </NativeSelectOption>
                    ))}
                  </NativeSelect>
                  <FieldDescription>{tr("应用模板时预填该分类，仍可编辑。", "This category is prefilled when the template is applied (still editable).")}</FieldDescription>
                </Field>
                <Field>
                  <FieldLabel>{tr("任务级拦截 / 放行规则", "Task-level intercept / allow rules")}</FieldLabel>
                  <AssetInterceptRulesEditor
                    value={draft.interceptRules}
                    onChange={(rules) => patchDraft({ interceptRules: rules })}
                  />
                  <FieldDescription>
                    {tr("应用模板时预填这些任务级规则（拦截/放行，仅对新任务生效，不影响全局）。", "These task-level rules are prefilled when the template is applied (intercept/allow, effective only for new tasks, never global).")}
                  </FieldDescription>
                </Field>
              </FieldGroup>
            </ScrollArea>
          </div>
          <SheetFooter className="border-t px-6 py-4 sm:flex-row sm:items-center">
            {selectedID != null && (
              <Button type="button" variant="destructive" className="sm:mr-auto" onClick={() => setDeleteOpen(true)}>
                <Trash2Icon data-icon="inline-start" />
                {tr("删除模板", "Delete template")}
              </Button>
            )}
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)}>
              {tr("关闭", "Close")}
            </Button>
            <Button type="button" disabled={saving} onClick={() => void save()}>
              {saving ? <Spinner data-icon="inline-start" /> : <SaveIcon data-icon="inline-start" />}
              {saveLabel}
            </Button>
          </SheetFooter>
        </SheetContent>
      </Sheet>
      <AlertDialog open={deleteOpen} onOpenChange={setDeleteOpen}>
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{tr(`删除模板「${draft.name || "未命名模板"}」？`, `Delete template "${draft.name || "Untitled template"}"?`)}</AlertDialogTitle>
            <AlertDialogDescription>{tr("已使用该模板创建的任务不受影响。", "Tasks already created from this template are not affected.")}</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel disabled={deleting}>{tr("取消", "Cancel")}</AlertDialogCancel>
            <AlertDialogAction
              variant="destructive"
              disabled={deleting}
              onClick={(event) => {
                event.preventDefault();
                void remove();
              }}
            >
              {deleting && <Spinner data-icon="inline-start" />}
              {deleting ? tr("正在删除", "Deleting") : tr("删除", "Delete")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>
    </>
  );
}

interface TaskTemplateControlsProps {
  description: string;
  goal: string;
  categoryID?: number;
  interceptRules?: AssetInterceptRuleInput[];
  selectedTemplateID: number | null;
  onSelectedTemplateIDChange: (id: number | null) => void;
  onApply: (template: TaskTemplate) => void;
  portalContainer?: React.RefObject<HTMLElement | null>;
}

export function TaskTemplateControls({
  description,
  goal,
  categoryID,
  interceptRules,
  selectedTemplateID,
  onSelectedTemplateIDChange,
  onApply,
  portalContainer,
}: TaskTemplateControlsProps) {
  const locale = useUiLocale();
  const tr = (zh: string, en: string) => uiText(locale, zh, en);
  const [templates, setTemplates] = React.useState<TaskTemplate[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [templateInputValue, setTemplateInputValue] = React.useState("");
  const [pendingTemplate, setPendingTemplate] = React.useState<TaskTemplate | null>(null);
  const [managerOpen, setManagerOpen] = React.useState(false);
  const [managerSeed, setManagerSeed] = React.useState<TemplateSeed | null>(null);

  const loadTemplates = React.useCallback(async () => {
    setLoading(true);
    try {
      setTemplates(await api.taskTemplates());
    } catch (error) {
      setTemplates([]);
      toast.error(uiText(locale, `加载模板失败：${(error as Error).message}`, `Failed to load templates: ${(error as Error).message}`));
    } finally {
      setLoading(false);
    }
  }, [locale]);

  React.useEffect(() => {
    void loadTemplates();
  }, [loadTemplates]);

  const selectedTemplate = React.useMemo(
    () => templates.find((template) => template.id === selectedTemplateID) ?? null,
    [selectedTemplateID, templates],
  );

  React.useEffect(() => {
    setTemplateInputValue(selectedTemplate?.name ?? "");
  }, [selectedTemplate]);

  const applyTemplate = (template: TaskTemplate) => {
    onApply(template);
    onSelectedTemplateIDChange(template.id);
    setPendingTemplate(null);
  };

  const chooseTemplate = (template: TaskTemplate | null) => {
    if (!template) {
      setTemplateInputValue("");
      onSelectedTemplateIDChange(null);
      return;
    }
    setTemplateInputValue(template.name);
    const hasContent = description.trim() !== "" || goal.trim() !== "";
    const changesContent = description !== template.description || goal !== template.goal;
    if (hasContent && changesContent) {
      setPendingTemplate(template);
      return;
    }
    applyTemplate(template);
  };

  const openManager = (seed: TemplateSeed | null) => {
    setManagerSeed(seed);
    setManagerOpen(true);
  };

  const upsertTemplate = (template: TaskTemplate) => {
    setTemplates((current) => [template, ...current.filter((item) => item.id !== template.id)]);
  };

  const deleteTemplate = (id: number) => {
    setTemplates((current) => current.filter((template) => template.id !== id));
    if (selectedTemplateID === id) onSelectedTemplateIDChange(null);
  };

  let pickerPlaceholder = loading ? tr("正在加载模板", "Loading templates") : tr("暂无任务模板", "No task templates yet");
  if (!loading && templates.length > 0) pickerPlaceholder = tr("搜索并选择任务模板", "Search and select a task template");

  return (
    <>
      <Field>
        <div className="flex flex-wrap items-center justify-between gap-2">
          <FieldLabel htmlFor="task-template-picker">{tr("任务模板", "Task template")}</FieldLabel>
          <div className="flex flex-wrap gap-2">
            <Button type="button" variant="ghost" size="sm" onClick={() => openManager(null)}>
              <Settings2Icon data-icon="inline-start" />
              {tr("管理模板", "Manage templates")}
            </Button>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              disabled={!description.trim() || !goal.trim()}
              onClick={() =>
                openManager({
                  description,
                  goal,
                  categoryID: categoryID ?? null,
                  interceptRules: interceptRules ?? [],
                })
              }
            >
              <SaveIcon data-icon="inline-start" />
              {tr("另存为模板", "Save as template")}
            </Button>
          </div>
        </div>
        <Combobox
          items={templates}
          itemToStringLabel={(template) => template.name}
          itemToStringValue={(template) => String(template.id)}
          inputValue={templateInputValue}
          onInputValueChange={(value) => setTemplateInputValue(value)}
          value={selectedTemplate}
          onValueChange={chooseTemplate}
        >
          <ComboboxInput
            id="task-template-picker"
            className="w-full"
            placeholder={pickerPlaceholder}
            disabled={loading || templates.length === 0}
            showClear
          />
          <ComboboxContent portalContainer={portalContainer}>
            <ComboboxEmpty>{tr("没有匹配的模板", "No matching templates")}</ComboboxEmpty>
            <ComboboxList>
              {(template) => (
                <ComboboxItem key={template.id} value={template}>
                  <LibraryIcon />
                  <div className="min-w-0 flex-1">
                    <p className="truncate font-medium text-sm">{template.name}</p>
                    {template.description && (
                      <p className="truncate text-muted-foreground text-xs">{template.description}</p>
                    )}
                  </div>
                </ComboboxItem>
              )}
            </ComboboxList>
          </ComboboxContent>
        </Combobox>
        <FieldDescription>{tr("选择后复制描述、目标、分类和任务级规则，不保留与模板的关联。", "Selecting copies the template's description, goal, category, and task-level rules; no link to the template is kept.")}</FieldDescription>
      </Field>

      <AlertDialog
        open={pendingTemplate != null}
        onOpenChange={(open) => {
          if (open) return;
          setPendingTemplate(null);
          setTemplateInputValue(selectedTemplate?.name ?? "");
        }}
      >
        <AlertDialogContent>
          <AlertDialogHeader>
            <AlertDialogTitle>{tr(`使用模板「${pendingTemplate?.name}」？`, `Use template "${pendingTemplate?.name}"?`)}</AlertDialogTitle>
            <AlertDialogDescription>{tr("已填写的描述和目标将被模板覆盖。", "The description and goal already filled in will be overwritten by the template.")}</AlertDialogDescription>
          </AlertDialogHeader>
          <AlertDialogFooter>
            <AlertDialogCancel>{tr("取消", "Cancel")}</AlertDialogCancel>
            <AlertDialogAction onClick={() => pendingTemplate && applyTemplate(pendingTemplate)}>
              {tr("覆盖并使用", "Overwrite and use")}
            </AlertDialogAction>
          </AlertDialogFooter>
        </AlertDialogContent>
      </AlertDialog>

      <TaskTemplateManager
        open={managerOpen}
        onOpenChange={setManagerOpen}
        templates={templates}
        seed={managerSeed}
        onCreated={upsertTemplate}
        onUpdated={upsertTemplate}
        onDeleted={deleteTemplate}
      />
    </>
  );
}
