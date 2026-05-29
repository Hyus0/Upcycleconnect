import { ref } from 'vue';

const t = ref({});
const currentLangCode = ref(localStorage.getItem("langCode") || "fr");
const API_URL = "http://localhost:8081";

const loadTranslations = async () => {
    try {
        const res = await fetch(`${API_URL}/traductions/${currentLangCode.value}`);
        if (res.ok) {
            t.value = await res.json();
        }
    } catch (e) {
        console.error("Erreur chargement traductions:", e);
    }
};

window.addEventListener("lang-changed", () => {
    currentLangCode.value = localStorage.getItem("langCode") || "fr";
    loadTranslations();
});

loadTranslations();

export function useTraduction() {
    return { 
        t, 
        currentLangCode,
        loadTranslations
    };
}