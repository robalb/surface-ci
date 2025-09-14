<template>
  <div class="asm-app">
    <!-- Left palette -->
    <aside class="sidebar">
      <h2>Nodes</h2>
      <div
        class="node-pill"
        draggable="true"
        @dragstart="onDragStart('find-subdomains', $event)"
        title="Input: Domain (text or upstream) → Output: list of subdomains"
      >
        Find Subdomains
      </div>

      <div
        class="node-pill"
        draggable="true"
        @dragstart="onDragStart('find-wildcard-domains', $event)"
        title="Input: Domain (text or upstream) → Output: wildcard domain patterns"
      >
        Find Wildcard Domains
      </div>

      <div
        class="node-pill"
        draggable="true"
        @dragstart="onDragStart('join', $event)"
        title="Input: two lists → Output: union(listA, listB)"
      >
        Join
      </div>
    </aside>

    <!-- Rete canvas -->
    <main
      ref="reteEl"
      class="canvas"
      @dragover.prevent
      @drop="onDrop"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount, ref } from 'vue'
import { NodeEditor } from 'rete'
import { AreaPlugin, AreaExtensions } from 'rete-area-plugin'
import { VuePlugin, VueArea2D, Presets as VuePresets } from 'rete-vue-plugin'
import { ConnectionPlugin, Presets as ConnectionPresets } from 'rete-connection-plugin'

import type { Schemes } from '../nodes'
import { createNode } from '../nodes'

type AreaExtra = VueArea2D<Schemes>

const reteEl = ref<HTMLElement | null>(null)

let editor: NodeEditor<Schemes> | null = null
let area: AreaPlugin<Schemes, AreaExtra> | null = null

function onDragStart(kind: string, e: DragEvent) {
  e.dataTransfer?.setData('application/x-rete-node-type', kind)
  e.dataTransfer?.setData('text/plain', kind) // fallback for some browsers
  e.dataTransfer!.effectAllowed = 'copy'
}

function toWorldCoordinates(e: DragEvent) {
  const rect = reteEl.value!.getBoundingClientRect()
  const { x, y, k } = area!.area.transform
  const worldX = (e.clientX - rect.left - x) / k
  const worldY = (e.clientY - rect.top - y) / k
  return { x: worldX, y: worldY }
}

async function onDrop(e: DragEvent) {
  if (!editor || !area) return
  const type = e.dataTransfer?.getData('application/x-rete-node-type')
            || e.dataTransfer?.getData('text/plain')
  if (!type) return

  const pos = toWorldCoordinates(e)
  const node = createNode(type as any)

  await editor.addNode(node)
  await area.translate(node.id, pos)
}

onMounted(async () => {
  if (!reteEl.value) return

  editor = new NodeEditor<Schemes>()
  area = new AreaPlugin<Schemes, AreaExtra>(reteEl.value)

  const render = new VuePlugin<Schemes, AreaExtra>()
  render.addPreset(VuePresets.classic.setup())

  const connection = new ConnectionPlugin<Schemes, AreaExtra>()
  connection.addPreset(ConnectionPresets.classic.setup())

  editor.use(area)
  area.use(connection)
  area.use(render)

  // Optional: helpful defaults
  AreaExtensions.selectableNodes(area, AreaExtensions.selector(), {
    accumulating: AreaExtensions.accumulateOnCtrl()
  })
  // Small hint grid
  reteEl.value.style.background = `
    radial-gradient(circle at 1px 1px, rgba(0,0,0,.15) 1px, transparent 0) 0 0 / 20px 20px
  `
})

onBeforeUnmount(() => {
  // Clean up
  area?.destroy()
  editor?.destroy()
})
</script>

<style scoped>
.asm-app {
  display: grid;
  grid-template-columns: 280px 1fr;
  grid-template-rows: 100%;
  height: 100vh;
  overflow: hidden;
}

.sidebar {
  border-right: 1px solid #e5e7eb;
  padding: 16px;
  background: #fafafa;
}

.sidebar h2 {
  margin: 0 0 12px 0;
  font: 600 14px/1.2 system-ui, -apple-system, Segoe UI, Roboto, sans-serif;
  letter-spacing: .02em;
  color: #111827;
}

.node-pill {
  user-select: none;
  cursor: grab;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  padding: 10px 12px;
  margin-bottom: 10px;
  background: #fff;
  font: 500 13px/1.1 system-ui, -apple-system, Segoe UI, Roboto, sans-serif;
}
.node-pill:active {
  cursor: grabbing;
}

.canvas {
  position: relative;
  height: 100%;
  /* Rete mounts its content inside */
}
</style>
