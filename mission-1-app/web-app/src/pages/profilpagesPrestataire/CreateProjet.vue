<template>
    <div class="page-container">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">ACCUEIL > MES PROJETS > CRÉER</p>
                <h1 class="hero-title1">DÉPOSER UN PROJET</h1>
                <p class="classic-text">
                    Partagez votre création upcyclée et inspirez la communauté.
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

        <form @submit.prevent="handleSubmit" class="create-annonce-form">
            <input
                type="file"
                ref="fileInput"
                @change="handleFileUpload"
                accept="image/*"
                style="display: none"
            />

            <div class="split-layout">
                <div class="form-card main-info-card">
                    <h2 class="card-title">1. Le Projet</h2>
                    <p class="card-subtitle">
                        Décrivez votre création et son impact.
                    </p>

                    <div class="form-group">
                        <label>Titre du projet</label>
                        <input
                            v-model="form.titre"
                            type="text"
                            required
                            placeholder="Ex: Table basse en palettes recyclées"
                        />
                    </div>

                    <div class="form-group">
                        <label>Description courte</label>
                        <textarea
                            v-model="form.description_courte"
                            rows="3"
                            maxlength="255"
                            required
                            placeholder="Résumez votre projet en quelques mots..."
                        ></textarea>
                        <small class="char-count"
                            >{{ form.description_courte.length }}/255</small
                        >
                    </div>

                    <div class="form-group">
                        <label>Image de couverture</label>
                        <div class="image-upload-wrapper">
                            <button
                                type="button"
                                class="btn-import"
                                @click="triggerUpload('couverture')"
                            >
                                Choisir une image
                            </button>
                            <img
                                v-if="form.image_url"
                                :src="form.image_url"
                                class="img-preview"
                                alt="Couverture"
                            />
                            <span v-else class="no-image-text"
                                >Aucune image sélectionnée</span
                            >
                        </div>
                    </div>

                    <div class="form-divider">Écologie & Visibilité</div>

                    <div class="form-row align-center">
                        <div class="form-group price-group">
                            <label>CO2 Évité estimé</label>
                            <div class="price-input-wrapper">
                                <input
                                    v-model.number="form.co2_evite_kg"
                                    type="number"
                                    step="0.1"
                                    min="0"
                                    placeholder="0.0"
                                />
                                <span class="currency-symbol">kg</span>
                            </div>
                        </div>

                        <div class="form-group">
                            <label>Visibilité</label>
                            <div class="radio-group">
                                <label class="radio-label">
                                    <input
                                        type="radio"
                                        v-model="form.visible_public"
                                        :value="true"
                                    />
                                    Public
                                </label>
                                <label class="radio-label">
                                    <input
                                        type="radio"
                                        v-model="form.visible_public"
                                        :value="false"
                                    />
                                    Privé
                                </label>
                            </div>
                        </div>
                        
                    </div>
                    <div class="form-divider">Mise en vente</div>
                    
                    <div class="form-group">
                        <label class="radio-label" style="margin-bottom: 0.5rem;">
                            <input type="checkbox" v-model="form.en_vente" style="width: 1.2rem; height: 1.2rem; accent-color: #2d7a4f;" />
                            Mettre en vente mon produit
                        </label>
                    
                        <div v-if="form.en_vente" class="price-input-wrapper" style="margin-top: 0.6rem;">
                            <input
                                v-model.number="form.prix"
                                type="number"
                                step="0.5"
                                min="0"
                                placeholder="Prix de vente"
                            />
                            <span class="currency-symbol">€</span>
                        </div>
                    </div>
                </div>

                <div class="right-column">
                    <div class="steps-global-header">
                        <div class="steps-header-text">
                            <h2 class="card-title">2. Étapes de fabrication</h2>
                            <p class="card-subtitle">
                                Ajoutez un tutoriel pas-à-pas (optionnel).
                            </p>
                        </div>
                        <button
                            type="button"
                            class="btn-add-step"
                            @click="addEtape"
                        >
                            + Ajouter
                        </button>
                    </div>

                    <div v-if="form.etapes.length === 0" class="no-steps-msg">
                        Aucune étape ajoutée.
                    </div>

                    <div class="steps-container">
                        <div
                            v-for="(etape, index) in form.etapes"
                            :key="index"
                            class="form-card step-full-card"
                        >
                            <div class="step-header-clean">
                                <div class="step-number-wrapper">
                                    <span class="step-number-badge">{{
                                        index + 1
                                    }}</span>
                                    <span class="step-title-label">Étape</span>
                                </div>
                                <button
                                    type="button"
                                    class="btn-delete-clean"
                                    @click="removeEtape(index)"
                                >
                                    Retirer
                                </button>
                            </div>

                            <div class="step-body-clean">
                                <div class="form-group">
                                    <label>Titre de l'étape</label>
                                    <input
                                        v-model="etape.titre"
                                        type="text"
                                        required
                                        placeholder="Titre de l'étape"
                                    />
                                </div>

                                <div class="form-group">
                                    <label>Description de l'étape</label>
                                    <textarea
                                        v-model="etape.description"
                                        rows="3"
                                        required
                                        placeholder="Détails des instructions..."
                                    ></textarea>
                                </div>

                                <div class="form-group">
                                    <label
                                        >Image de l'étape (optionnelle)</label
                                    >
                                    <div class="image-upload-wrapper">
                                        <button
                                            type="button"
                                            class="btn-import"
                                            @click="
                                                triggerUpload('etape', index)
                                            "
                                        >
                                            Ajouter une image
                                        </button>
                                        <img
                                            v-if="etape.image_url"
                                            :src="etape.image_url"
                                            class="img-preview"
                                            alt="Étape"
                                        />
                                    </div>
                                </div>
                            </div>
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
                            {{
                                loading
                                    ? "Publication..."
                                    : "Publier le projet"
                            }}
                        </button>
                    </div>
                </div>
            </div>
        </form>
    </div>
</template>

<script setup>
import { ref, computed } from "vue";
import { useRouter } from "vue-router";

const router = useRouter();
const loading = ref(false);
const errors = ref([]);
const successMsg = ref("");

const API_URL = "/go";

const currentUserId = computed(() => {
    const storedId =
        sessionStorage.getItem("id") || sessionStorage.getItem("userId");
    return Number(storedId) || 0;
});

const form = ref({
    titre: "",
    description_courte: "",
    image_url: "",
    co2_evite_kg: 0,
    visible_public: true,
    en_vente: false,     
    prix: 0,
    etapes: [],
});

const fileInput = ref(null);
const uploadContext = ref({ type: "", index: null });

const triggerUpload = (type, index = null) => {
    uploadContext.value = { type, index };
    fileInput.value.click();
};

const handleFileUpload = async (event) => {
    const file = event.target.files[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("image", file);

    try {
        const token = sessionStorage.getItem("userToken");

        const res = await fetch(`${API_URL}/projets/upload-image`, {
            method: "POST",
            headers: {
                Authorization: token,
            },
            body: formData,
        });

        if (res.ok) {
            const data = await res.json();
            const { type, index } = uploadContext.value;

            if (type === "couverture") {
                form.value.image_url = data.url;
            } else if (type === "etape" && index !== null) {
                form.value.etapes[index].image_url = data.url;
            }
        } else {
            alert("Erreur lors de l'upload de l'image.");
        }
    } catch (e) {
        console.error("Erreur d'upload :", e);
        alert("Serveur injoignable.");
    } finally {
        event.target.value = "";
    }
};

const addEtape = () => {
    form.value.etapes.push({ titre: "", description: "", image_url: "" });
};

const removeEtape = (index) => {
    form.value.etapes.splice(index, 1);
};

const validateFrontend = () => {
    errors.value = [];
    if (form.value.titre.length < 5)
        errors.value.push("Le titre doit faire au moins 5 caractères.");
    if (!form.value.image_url)
        errors.value.push("L'image de couverture est obligatoire.");
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
        id_createur: currentUserId.value,
        titre: form.value.titre,
        description_courte: form.value.description_courte,
        image_url: form.value.image_url,
        co2_evite_kg: parseFloat(form.value.co2_evite_kg) || 0,
        visible_public: form.value.visible_public,
        score_impact: (parseFloat(form.value.co2_evite_kg) || 0) * 10,
        prix: form.value.en_vente ? parseFloat(form.value.prix) || 0 : null,
        etapes: form.value.etapes.map((etape, index) => ({
            numero_ordre: index + 1,
            titre: etape.titre,
            description: etape.description,
            image_url: etape.image_url || "",
        })),
    };

    try {
        const response = await fetch(`${API_URL}/projets`, {
            method: "POST",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(payload),
        });

        if (response.ok) {
            successMsg.value = "Projet publié avec succès !";
            setTimeout(() => {
                router.push("/profil/projets");
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
    margin: 0.2rem 0 1.5rem 0;
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

.form-divider {
    margin-top: 1rem;
    padding-bottom: 0.5rem;
    border-bottom: 2px solid #e9f5ed;
    font-weight: 800;
    color: #2d7a4f;
    font-size: 0.85rem;
    text-transform: uppercase;
    letter-spacing: 1px;
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

.image-upload-wrapper {
    display: flex;
    align-items: center;
    gap: 15px;
}
.btn-import {
    background: #f5f5f5;
    border: 1px dashed #ccc;
    padding: 0.8rem 1.2rem;
    border-radius: 10px;
    font-weight: 600;
    cursor: pointer;
    color: #555;
    transition: all 0.2s;
}
.btn-import:hover {
    background: #e0e0e0;
    border-color: #aaa;
}
.img-preview {
    width: 120px;
    height: 80px;
    object-fit: cover;
    border-radius: 8px;
    border: 1px solid #ddd;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}
.no-image-text {
    font-size: 0.85rem;
    color: #999;
    font-style: italic;
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

.price-input-wrapper {
    position: relative;
    display: flex;
    align-items: center;
}
.price-input-wrapper input {
    width: 100%;
    padding-right: 2.5rem;
}
.currency-symbol {
    position: absolute;
    right: 1rem;
    font-weight: bold;
    color: #888;
}

.steps-global-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1rem;
}
.btn-add-step {
    background: none;
    border: 1px solid #2d7a4f;
    color: #2d7a4f;
    padding: 6px 12px;
    border-radius: 8px;
    font-weight: bold;
    cursor: pointer;
    font-size: 0.85rem;
}
.btn-add-step:hover {
    background: #e9f5ed;
}
.no-steps-msg {
    color: #888;
    font-size: 0.9rem;
    font-style: italic;
}

.steps-container {
    display: flex;
    flex-direction: column;
    gap: 0;
    width: 100%;
    box-sizing: border-box;
}

.step-full-card {
    padding: 1.5rem;
    gap: 1rem;
    margin-bottom: 1.5rem;
    width: 100%;
    box-sizing: border-box;
}
.step-header-clean {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-bottom: 1rem;
    border-bottom: 1px dashed #eee;
    margin-bottom: 0.5rem;
}
.step-number-badge {
    background: #2d7a4f;
    color: white;
    font-weight: 700;
    font-size: 0.85rem;
    width: 26px;
    height: 26px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
}
.step-title-label {
    font-size: 1.1rem;
    font-weight: 700;
    color: #1a1a1a;
    margin: 0;
}
.btn-delete-clean {
    background: none;
    border: none;
    color: #ef4444;
    cursor: pointer;
    font-size: 0.85rem;
    font-weight: 600;
}
.btn-delete-clean:hover {
    text-decoration: underline;
}
.step-body-clean {
    display: flex;
    flex-direction: column;
    gap: 1rem;
}

.form-actions-card {
    margin-top: 1rem;
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
.btn-save:hover {
    background-color: #246343;
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
.char-count {
    text-align: right;
    color: #999;
    font-size: 0.8rem;
    margin-top: -5px;
}
</style>
