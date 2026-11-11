"use client";

import * as React from "react";
import { AlertCircleIcon, RefreshCwIcon } from "lucide-react";

import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { ExplorationGraph } from "@/components/exploration-graph";
import { api } from "@/lib/api";
import type { Edge, TaskNode } from "@/lib/types";

export function GraphTab({ taskId }: { taskId: string }) {
  const [nodes, setNodes] = React.useState<TaskNode[]>([]);
  const [edges, setEdges] = React.useState<Edge[]>([]);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState("");
  const [refresh, setRefresh] = React.useState(0);
  // Signature of the previous graph data: polling skips setState on identical data to avoid
  // rebuilding the whole graph (so dragging is not jarred by the 20s poll). Only render-affecting fields count.
  const sigRef = React.useRef("");

  React.useEffect(() => {
    let cancelled = false;
    sigRef.current = "";
    setNodes([]);
    setEdges([]);
    setError("");
    setLoading(true);

    const load = (showProgress: boolean) => {
      if (showProgress) setLoading(true);
      api
        .explorationGraph(taskId)
        .then((g) => {
          if (cancelled) return;
          setError("");
          const ns = g.nodes ?? [];
          const es = g.edges ?? [];
          const sig = JSON.stringify([
            ns.map((n) => [n.id, n.type, n.state, n.priority, n.payload]),
            es.map((e) => [e.src, e.dst, e.rel]),
          ]);
          if (sig === sigRef.current) return; // avoid re-layout while dragging
          sigRef.current = sig;
          setNodes(ns);
          setEdges(es);
        })
        .catch((cause) => {
          if (!cancelled) setError(cause instanceof Error ? cause.message : "图谱接口请求失败");
          // Keep the last successful graph visible instead of turning it into an empty canvas.
        })
        .finally(() => {
          if (!cancelled && showProgress) setLoading(false);
        });
    };
    load(true);
    const timer = setInterval(() => load(false), 20000);
    return (
    <Card className="min-w-0 gap-0 overflow-hidden rounded-xl border-border/70 bg-card/90 py-0 shadow-sm">
      <CardHeader className="flex flex-wrap items-center justify-between gap-3 border-b border-border/70 px-4 py-4 sm:px-5">
        <div className="min-w-0">
          <CardTitle className="text-base font-semibold">任务探索图谱</CardTitle>
          <p className="mt-1 text-xs text-muted-foreground">
            {loading ? "正在加载图谱…" : `${nodes.length} 个节点 · ${edges.length} 条关系 · 每 20 秒自动更新`}
          </p>
        </div>
        <Button size="sm" variant="outline" disabled={loading} onClick={() => setRefresh((n) => n + 1)}>
          <RefreshCwIcon className={loading ? "animate-spin" : ""} /> 刷新图谱
        </Button>
      </CardHeader>
      <CardContent className="min-w-0 p-0">
        {error && (
          <div role="alert" className="mx-4 mt-3 flex min-w-0 items-start gap-2 rounded-lg border border-amber-500/30 bg-amber-500/10 p-3 text-sm sm:mx-5">
            <AlertCircleIcon className="mt-0.5 size-4 shrink-0 text-amber-500" />
            <div className="min-w-0">
              <p className="font-medium">图谱更新失败，已保留上次成功加载的内容</p>
              <p className="mt-1 break-words text-muted-foreground">{error}</p>
            </div>
          </div>
        )}
        <ExplorationGraph
          nodes={nodes}
          edges={edges}
          className="h-[clamp(22rem,65dvh,52rem)] min-w-0"
          emptyHint={loading ? "图谱加载中…" : error ? "暂时无法加载图谱，请重试" : "该任务暂时没有探索图谱节点"}
        />
      </CardContent>
    </Card>
  );
}
