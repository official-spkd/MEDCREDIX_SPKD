// Toast system — padanan hooks/use-toast + Toaster.
import { defineStore } from 'pinia'

export type ToastItem = {
  id: number
  title: string
  description?: string
  variant?: 'default' | 'destructive'
}

export const useToastStore = defineStore('toast', {
  state: () => ({ items: [] as ToastItem[], seq: 0 }),
  actions: {
    push(t: Omit<ToastItem, 'id'>) {
      const id = ++this.seq
      this.items.push({ id, ...t })
      setTimeout(() => this.dismiss(id), 5000)
    },
    dismiss(id: number) {
      this.items = this.items.filter((i) => i.id !== id)
    },
  },
})

export function toast(t: { title: string; description?: string; variant?: 'default' | 'destructive' }) {
  useToastStore().push(t)
}
