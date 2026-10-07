// useIntel — ambil data intelligence (dashboard) untuk view kesiapan.
import { onMounted, ref, watch } from 'vue'
import { api } from '~/composables/useApi'
import { useMedCredix } from '~/stores/medcredix'
export function useIntel() {
  const store = useMedCredix()
  const data = ref<any>(null)
  const loading = ref(true)
  async function load() {
    loading.value = true
    try {
      data.value = await api<any>(`/dashboard${store.hospitalId ? `?hospitalId=${store.hospitalId}` : ''}`)
    } catch { data.value = null }
    loading.value = false
  }
  onMounted(load)
  watch(() => store.hospitalId, load)
  return { data, loading, reload: load }
}
