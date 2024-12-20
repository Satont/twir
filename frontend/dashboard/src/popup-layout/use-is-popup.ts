import { createGlobalState } from '@vueuse/core'
import { readonly, ref } from 'vue'

export const useIsPopup = createGlobalState(() => {
	const isPopup = ref(false)

	function setIsPopup(v: boolean) {
		isPopup.value = v
	}

	return {
		isPopup: readonly(isPopup),
		setIsPopup,
	}
})
