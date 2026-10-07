<script setup lang="ts">
// Toaster — wadah notifikasi toast (pojok kanan bawah).
import { CheckCircle2, XCircle, X } from 'lucide-vue-next'
import { useToastStore } from '~/composables/useToast'
const store = useToastStore()
</script>

<template>
  <Teleport to="body">
    <div class="fixed bottom-4 right-4 z-[100] flex w-full max-w-sm flex-col gap-2">
      <TransitionGroup name="toast">
        <div
          v-for="t in store.items"
          :key="t.id"
          :class="[
            'pointer-events-auto flex items-start gap-2.5 rounded-xl border p-3.5 shadow-lg backdrop-blur',
            t.variant === 'destructive' ? 'border-danger/40 bg-danger-light text-danger' : 'border-border bg-card text-foreground',
          ]"
          role="status"
        >
          <CheckCircle2 v-if="t.variant !== 'destructive'" class="mt-0.5 h-4 w-4 shrink-0 text-success" />
          <XCircle v-else class="mt-0.5 h-4 w-4 shrink-0" />
          <div class="min-w-0 flex-1">
            <div class="text-xs font-semibold leading-snug">{{ t.title }}</div>
            <div v-if="t.description" class="mt-0.5 text-[11px] leading-relaxed opacity-80">{{ t.description }}</div>
          </div>
          <button class="shrink-0 opacity-50 hover:opacity-100" aria-label="Tutup" @click="store.dismiss(t.id)">
            <X class="h-3.5 w-3.5" />
          </button>
        </div>
      </TransitionGroup>
    </div>
  </Teleport>
</template>

<style scoped>
.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}
.toast-enter-from,
.toast-leave-to {
  opacity: 0;
  transform: translateY(8px);
}
</style>
