<template>
    <div class="page-container">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > TIPS > MODIFIER</p>
                <h1 class="hero-title1">MODIFIER UN TIP</h1>
                <p class="classic-text">
                    Mettez à jour les informations de votre astuce.
                </p>

                <div
                    v-if="errors.length > 0"
                    class="error-box"
                    style="margin-top: 1rem; margin-bottom: 0"
                >
                    <ul style="margin: 0; padding-left: 20px">
                        <li v-for="(err, index) in errors" :key="index">
                            {{ err }}
                        </li>
                    </ul>
                </div>

                <div
                    v-if="successMsg"
                    class="success-box"
                    style="margin-top: 1rem; margin-bottom: 0"
                >
                    {{ successMsg }}
                </div>
            </div>

            <button class="btn-secondary" @click="$router.back()">
                🠔 Retour
            </button>
        </header>

        <div v-if="isFetching" class="loading-state" style="margin-top: 2rem;">
            Chargement des données du tip...
        </div>

        <form v-else @submit.prevent="handleSubmit" class="create-annonce-form">
            <div class="split-layout">
                <div class="form-card main-info-card">
                    <h2 class="card-title">1. Le Tip</h2>
                    <p class="card-subtitle">
                        Décrivez votre astuce et son public cible.
                    </p>

                    <div class="form-group">
                        <label>Titre du tip</label>
                        <input
                            v-model="form.titre"
                            type="text"
                            required
                            placeholder="Ex: Comment trier ses déchets efficacement"
                        />
                    </div>

                    <div class="form-group">
                        <label>Description</label>
                        <textarea
                            v-model="form.description"
                            rows="6"
                            required
                            placeholder="Expliquez votre astuce en détail..."
                        ></textarea>
                    </div>

                    <div class="form-divider">Diffusion & Visibilité</div>

                    <div class="form-row">
                        <div class="form-group">
                            <label>Public cible</label>
                            <select v-model="form.role_cible" required>
                                <option value="" disabled>Sélectionner...</option>
                                <option value="Particulier">Particulier</option>
                                <option value="Prestataire">Prestataire</option>
                                <option value="Salarie">Salarié</option>
                                <option value="Admin">Admin</option>
                            </select>
                        </div>

                        <div class="form-group">
                            <label>Statut</label>
                            <div class="radio-group">
                                <label class="radio-label">
                                    <input
                                        type="radio"
                                        v-model="form.actif"
                                        :value="true"
                                    />
                                    Actif
                                </label>
                                <label class="radio-label">
                                    <input
                                        type="radio"
                                        v-model="form.actif"
                                        :value="false"
                                    />
                                    Inactif
                                </label>
                            </div>
                        </div>
                    </div>
                </div>

                <div class="right-column">
                    <div class="form-card">
                        <h2 class="card-title">2. Contenu multimédia</h2>
                        <p class="card-subtitle">
                            Ajoutez une vidéo pour illustrer votre tip (optionnel).
                        </p>

                        <div class="form-group">
                            <label>URL de la vidéo</label>
                            <input
                                v-model="form.video_url"
                                type="url"
                                placeholder="https://youtube.com/watch?v=..."
                            />
                        </div>

                        <div v-if="videoEmbedUrl" class="video-preview-wrapper">
                            <label class="info-label">Aperçu</label>
                            <iframe
                                :src="videoEmbedUrl"
                                class="video-preview"
                                frameborder="0"
                                allowfullscreen
                            ></iframe>
                        </div>

                        <div v-else-if="form.video_url" class="video-fallback">
                            <span class="no-image-text">⚠️ URL non prévisualisable — elle sera enregistrée.</span>
                        </div>

                        <div v-else class="video-placeholder">
                            <span class="no-image-text">Aucune vidéo renseignée.</span>
                        </div>
                    </div>

                    <div class="form-actions-card">
                        <button
                            type="button"
                            class="btn-cancel"
                            @click="$router.back()"
                        >
                            Annuler
                        </button>
                        <button
                            type="submit"
                            class="btn-save"
                            :disabled="loading"
                        >
                            {{ loading ? "Enregistrement..." : "Enregistrer les modifications" }}
                        </button>
                    </div>
                </div>
            </div>
        </form>
    </div>
</template>

<script setup>
import { ref, computed, onMounted } from "vue";
import { useRouter, useRoute } from "vue-router";

const router = useRouter();
const route = useRoute(); 
const loading = ref(false);
const isFetching = ref(true);
const errors = ref([]);
const successMsg = ref("");

const API_URL = "http://localhost:8081";
const tipId = route.params.id;

const currentUserId = computed(() => {
    const storedId = sessionStorage.getItem("id") || sessionStorage.getItem("userId");
    return Number(storedId) || 0;
});

const form = ref({
    titre: "",
    description: "",
    video_url: "",
    role_cible: "",
    actif: true,
});

const fetchTipData = async () => {
    const token = sessionStorage.getItem("userToken");
    
    try {
        const res = await fetch(`${API_URL}/tips/${tipId}`, {
            method: "GET",
            headers: { Authorization: token },
        });

        if (res.ok) {
            const data = await res.json();
            form.value.titre = data.titre || "";
            form.value.description = data.description || "";
            form.value.video_url = data.video_url || "";
            form.value.role_cible = data.role_cible || "";
            form.value.actif = !!data.actif; 
        } else {
            errors.value.push("Impossible de charger les informations de ce tip.");
        }
    } catch (err) {
        console.error("Erreur de chargement:", err);
        errors.value.push("Erreur de connexion au serveur.");
    } finally {
        isFetching.value = false;
    }
};

onMounted(() => {
    fetchTipData();
});

const videoEmbedUrl = computed(() => {
    const url = form.value.video_url;
    if (!url) return null;
    const ytMatch = url.match(
        /(?:youtube\.com\/watch\?v=|youtu\.be\/)([a-zA-Z0-9_-]{11})/
    );
    if (ytMatch) return `https://www.youtube.com/embed/${ytMatch[1]}`;
    return null;
});

const validateFrontend = () => {
    errors.value = [];
    if (form.value.titre.length < 5)
        errors.value.push("Le titre doit faire au moins 5 caractères.");
    if (form.value.description.length < 10)
        errors.value.push("La description doit faire au moins 10 caractères.");
    if (!form.value.role_cible)
        errors.value.push("Veuillez sélectionner un public cible.");
    return errors.value.length === 0;
};

const handleSubmit = async () => {
    if (currentUserId.value === 0) {
        errors.value = ["Erreur de session. Reconnectez-vous."];
        return;
    }

    if (!validateFrontend()) return;

    loading.value = true;
    successMsg.value = "";
    errors.value = [];

    const token = sessionStorage.getItem("userToken");

    const payload = {
        id: Number(tipId),
        id_createur: currentUserId.value,
        titre: form.value.titre,
        description: form.value.description,
        video_url: form.value.video_url || "",
        role_cible: form.value.role_cible,
        actif: form.value.actif 
    };

    try {
        const response = await fetch(`${API_URL}/tips/${tipId}`, {
            method: "PUT",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });

        if (response.ok) {
            successMsg.value = "Tip modifié avec succès !";
            setTimeout(() => {
                router.push("/profil/tips"); 
            }, 1500);
        } else {
            const errMsg = await response.text();
            errors.value.push("Erreur serveur : " + errMsg);
        }
    } catch (err) {
        console.error("Erreur:", err);
        errors.value.push("Impossible de joindre le serveur.");
    } finally {
        loading.value = false;
    }
};
</script>

<style scoped>
.page-container {
    padding: 20px;
}

.content-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 2rem;
}

.header-left {
    display: flex;
    flex-direction: column;
}

.sidebar__category2 {
    font-size: 0.65rem;
    color: #8fa396;
    letter-spacing: 1px;
    margin: 0 0 0.5rem 0;
    text-transform: uppercase;
}

.hero-title1 {
    font-size: 2rem;
    font-weight: 800;
    margin: 1.5rem 0 0.5rem;
    color: #1a1a1a;
}

.classic-text {
    color: #666;
    margin: 0;
}

.btn-secondary {
    padding: 8px 16px;
    border-radius: 10px;
    border: 1px solid #ddd;
    background: white;
    cursor: pointer;
    font-weight: 500;
}

.error-box {
    background-color: #fee2e2;
    border: 1px solid #ef4444;
    color: #b91c1c;
    padding: 12px;
    border-radius: 10px;
    margin-bottom: 15px;
}

.success-box {
    background-color: #f0fdf4;
    border: 1px solid #22c55e;
    color: #166534;
    padding: 12px;
    border-radius: 10px;
    margin-bottom: 15px;
}

.create-annonce-form {
    width: 100%;
    box-sizing: border-box;
}

.split-layout {
    display: grid;
    grid-template-columns: 1.5fr 1fr;
    gap: 2rem;
    width: 100%;
}

.right-column {
    display: flex;
    flex-direction: column;
}

.form-card {
    background: #ffffff;
    padding: 2rem;
    border-radius: 16px;
    border: 1px solid #eee;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.03);
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    margin-bottom: 2rem;
    width: 100%;
    box-sizing: border-box;
}

.card-title {
    font-size: 1.4rem;
    font-weight: 700;
    margin: 0;
    color: #1a1a1a;
}

.card-subtitle {
    font-size: 0.9rem;
    color: #666;
    margin: 0.2rem 0 0 0;
}

.form-divider {
    margin-top: 0.5rem;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid #e9f5ed;
    font-weight: 800;
    color: #2d7a4f;
    font-size: 0.85rem;
    text-transform: uppercase;
    letter-spacing: 1px;
}

.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
}

.form-group {
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
}

.form-group label {
    font-size: 0.95rem;
    font-weight: 600;
    color: #333;
}

input,
select,
textarea {
    padding: 0.9rem;
    border: 1px solid #ddd;
    border-radius: 10px;
    font-family: inherit;
    font-size: 0.95rem;
    width: 100%;
    box-sizing: border-box;
}

input:focus,
select:focus,
textarea:focus {
    outline: none;
    border-color: #2d7a4f;
}

textarea {
    resize: vertical;
}

.radio-group {
    display: flex;
    gap: 1.5rem;
    padding: 0.5rem 0;
}

.radio-label {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    cursor: pointer;
    font-weight: 600;
}

.radio-label input[type="radio"] {
    width: 1.2rem;
    height: 1.2rem;
    accent-color: #2d7a4f;
}

.info-label {
    font-size: 0.85rem;
    color: #2d7a4f;
    text-transform: uppercase;
    font-weight: 800;
    letter-spacing: 0.5px;
}

.video-preview-wrapper {
    display: flex;
    flex-direction: column;
    gap: 0.6rem;
}

.video-preview {
    width: 100%;
    aspect-ratio: 16 / 9;
    border-radius: 10px;
    border: 1px solid #eee;
}

.video-placeholder,
.video-fallback {
    background: #fcfcfc;
    border: 1px dashed #ddd;
    border-radius: 10px;
    padding: 1.5rem;
    text-align: center;
}

.no-image-text {
    font-size: 0.85rem;
    color: #999;
    font-style: italic;
}

.recap-card {
    gap: 1rem;
}

.recap-row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    font-size: 0.9rem;
    padding: 0.4rem 0;
    border-bottom: 1px solid #f5f5f5;
}

.recap-row:last-child {
    border-bottom: none;
}

.recap-label {
    color: #888;
    font-weight: 500;
}

.recap-value {
    font-weight: 700;
    color: #333;
}

.text-success {
    color: #2d7a4f;
    font-weight: 700;
}

.type-badge {
    padding: 4px 10px;
    border-radius: 12px;
    font-size: 0.75rem;
    font-weight: 800;
}

.role-particulier  { background: #e9f5ed; color: #2d7a4f; }
.role-prestataire  { background: #e8eaf6; color: #3f51b5; }
.role-salarie      { background: #fff3e0; color: #e65100; }
.role-admin        { background: #fce4ec; color: #c2185b; }

.form-actions-card {
    margin-top: 0.5rem;
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
}

.btn-save {
    background-color: #2d7a4f;
    color: white;
    padding: 1rem 2.5rem;
    border: none;
    border-radius: 12px;
    font-weight: 700;
    cursor: pointer;
    transition: background 0.2s;
}

.btn-save:hover:not(:disabled) {
    background-color: #246343;
}

.btn-save:disabled {
    opacity: 0.6;
    cursor: not-allowed;
}

.btn-cancel {
    background: #ffffff;
    border: 1px solid #ccc;
    color: #666;
    padding: 1rem 2.5rem;
    border-radius: 12px;
    cursor: pointer;
    font-weight: 600;
}

.btn-cancel:hover {
    background-color: #f5f5f5;
    border-color: #999;
}

.loading-state {
    font-style: italic;
    color: #8fa396;
}
</style>