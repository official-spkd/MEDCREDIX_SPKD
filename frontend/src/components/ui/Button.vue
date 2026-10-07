<script setup lang="ts">
// Button — padanan shadcn/ui Button (New York style) dengan varian SPKD.
const props = defineProps<{
  variant?: 'default' | 'outline' | 'ghost' | 'secondary' | 'link' | 'brand'
  size?: 'default' | 'sm' | 'icon' | 'icon-sm'
  disabled?: boolean
  type?: 'button' | 'submit'
}>()
const variant = props.variant ?? 'default'
const size = props.size ?? 'default'
const variantCls: Record<string, string> = {
  default: 'bg-brand text-white shadow-xs hover:bg-brand-hover',
  outline: 'border border-input bg-card shadow-xs hover:bg-muted/60',
  ghost: 'hover:bg-muted/70',
  secondary: 'bg-muted text-foreground shadow-xs hover:bg-muted/70',
  link: 'text-brand-2 underline-offset-4 hover:underline',
  brand: 'bg-brand-2 text-white shadow-xs hover:bg-brand-2-hover',
}
const sizeCls: Record<string, string> = {
  default: 'h-9 px-4 py-2 text-sm',
  sm: 'h-8 gap-1.5 rounded-md px-3 text-xs',
  icon: 'h-9 w-9',
  'icon-sm': 'h-8 w-8',
}
</script>

<template>
  <button
    :type="type ?? 'button'"
    :disabled="disabled"
    :class="[
      'inline-flex items-center justify-center gap-2 whitespace-nowrap rounded-lg text-sm font-semibold transition-colors',
      'focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring disabled:pointer-events-none disabled:opacity-50',
      '[&_svg]:shrink-0',
      variantCls[variant],
      sizeCls[size],
    ]"
  >
    <slot />
  </button>
</template>
