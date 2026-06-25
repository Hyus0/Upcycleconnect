<template>
    <div class="info-view">
        <header class="content-header">
            <div class="header-left">
                <p class="sidebar__category2">
                    ACCUEIL > PARAMÈTRES > INFORMATIONS
                </p>
                <h1 class="hero-title1">MON PROFIL</h1>
                <p class="classic-text">
                    Gérez vos informations personnelles et les détails de votre
                    compte utilisateur.
                </p>
                <div v-if="errors.length > 0" 
                style="background-color: #fee2e2; border: 1px solid #ef4444; color: #b91c1c; 
                padding: 10px; border-radius: 8px; margin-bottom: 15px;">
                    <ul style="margin: 0; padding-left: 20px;">
                        <li v-for="(err, index) in errors" :key="index">
                            {{ err }}
                        </li>
                    </ul>
                </div>
                <div v-if="successMsg" class="success-box"
                style="background-color: #E2FEE3; border: 1px solid #44EF44; color: #158F3C; 
                padding: 10px; border-radius: 8px; margin-bottom: 15px;">
                    {{ successMsg }}
                </div>
            </div>
            <button class="btn-main-action" @click="updateProfile">
                Mettre à jour le profil
            </button>
        </header>

        <div class="info-layout">
            <div class="info-main-col">
                
                <section class="info-section">
                    <h2 class="section-title">Images du profil</h2>
                    
                    <div class="images-row">
                        <div class="upload-group avatar-group">
                            <label class="image-label">Photo de profil</label>
                            <div class="upload-controls">
                                <div class="avatar-container">
                                    <img :src="profilePreview || defaultAvatar" alt="Avatar" class="avatar-preview" />
                                </div>
                                <label for="profile-upload" class="btn-import">Importer</label>
                                <input 
                                    id="profile-upload" 
                                    type="file" 
                                    accept="image/png, image/jpeg, image/jpg" 
                                    @change="handleProfileUpload" 
                                    hidden 
                                />
                            </div>
                        </div>
                
                        <div class="upload-group banner-group">
                            <label class="image-label">Bannière</label>
                            <div class="upload-controls">
                                <div 
                                    class="image-preview banner-preview" 
                                    :style="{ backgroundImage: 'url(' + (bannerPreview || defaultBanner) + ')' }"
                                >
                                </div>
                                <label for="banner-upload" class="btn-import">Importer</label>
                                <input 
                                    id="banner-upload" 
                                    type="file" 
                                    accept="image/png, image/jpeg, image/jpg" 
                                    @change="handleBannerUpload" 
                                    hidden 
                                />
                            </div>
                        </div>
                    </div>
                </section>

                <section class="info-section">
                    <h2 class="section-title">Identité Personnelle</h2>
                    
                    <div class="input-grid">
                        <div class="input-field">
                            <label>Prénom</label>
                            <input type="text" v-model="form.prenom" />
                        </div>
                        <div class="input-field">
                            <label>Nom</label>
                            <input type="text" v-model="form.nom" />
                        </div>
                        <div class="input-field">
                            <label>Date de naissance</label>
                            <input type="date" v-model="form.date_naissance" disabled class="disabled-input"/>
                        </div>
                        <div class="input-field">
                            <label class="siret-label">
                                Email (Identifiant)
                                <span v-if="form.mail_valide === true || form.mail_valide === 'true' || form.mail_valide === 1" class="badge-valid">Validé</span>
                                <span v-else class="badge-pending">En attente de validation</span>
                            </label>
                            <input 
                                type="email" 
                                v-model="form.mail" 
                            />
                        </div>
                    </div>

                    <div v-if="form.role === 'Prestataire'" class="input-field siret-box">
                        <label class="siret-label">
                            Numéro SIRET (14 chiffres)
                            <span v-if="form.siret_valide" class="badge-valid">Vérifié</span>
                            <span v-else class="badge-pending">En attente de vérification</span>
                        </label>
                        <input 
                            type="text" 
                            v-model="form.siret" 
                            maxlength="14" 
                            placeholder="Ex: 12345678901234" 
                            :disabled="form.siret_valide"
                            :class="{ 'disabled-input': form.siret_valide }"
                        />
                        <small class="warning-text" style="text-align: left; margin-top: 4px;">
                            Toute modification de votre SIRET nécessitera une nouvelle vérification par nos équipes.
                        </small>
                    </div>

                </section>
            
                <section class="info-section">
                    <h2 class="section-title">Adresse et Localisation</h2>
                    <div class="input-field">
                        <label>Adresse complète</label>
                        <input
                            type="text"
                            v-model="form.adresse"
                            placeholder="Ex: 15 rue des Artisans"
                        />
                    </div>
                    <div class="input-grid">
                        <div class="input-field">
                            <label>Ville</label>
                            <input type="text" v-model="form.ville" />
                        </div>
                        <div class="input-field">
                            <label>Code Postal</label>
                            <input type="text" v-model="form.code_postal" />
                        </div>
                    </div>
                </section>
            </div>

            <aside class="info-side-col">
                <section class="info-section status-card">
                    <h2 class="section-title">État du compte</h2>
                    <div class="status-item">
                        <span class="status-label">Rôle actuel</span>
                        <span class="badge-role">{{ form.role }}</span>
                    </div>
                    <div class="status-item">
                        <span class="status-label">Statut</span>
                        <span class="status-active">{{ form.statut }}</span>
                    </div>
                    <div class="status-item">
                        <span class="status-label">Langue ID</span>
                        <span>{{ form.id_langue }}</span>
                    </div>
                    <div class="divider"></div>
                    <div class="registration-info">
                        <p>Membre depuis le :</p>
                        <strong>{{ formatDate(form.date_inscription) }}</strong>
                    </div>
                </section>

                <section class="info-section danger-zone">
                    <h2 class="section-title">Sécurité</h2>
                    <button class="btn-outline" @click="router.push('/profil/password')">Changer le mot de passe</button>
                    <p class="warning-text">
                        Le mot de passe doit être modifié régulièrement.
                    </p>
                </section>
            </aside>
        </div>
    </div>
</template>

<script setup>
import { ref, onMounted } from "vue";
import { useRouter } from 'vue-router';
import basicAvatar from '../../components/basicAvatar.png';
import basicBanner from '../../components/basicBanner.jpg';

const router = useRouter();

const form = ref({
    prenom: "",
    nom: "",
    mail: "",
    mail_valide: false,
    adresse: "",
    ville: "",
    code_postal: "",
    date_naissance: "",
    date_inscription: "",
    role: "",
    statut: "",
    id_langue: 1,
    siret: "",            
    siret_valide: false,
});

const errors = ref([]);
const successMsg = ref("");

const profileFile = ref(null);
const bannerFile = ref(null);
const profilePreview = ref(null);
const bannerPreview = ref(null);

const defaultAvatar = basicAvatar; 
const defaultBanner = basicBanner;

const handleProfileUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
        profileFile.value = file;
        profilePreview.value = URL.createObjectURL(file);
    }
};

const handleBannerUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
        bannerFile.value = file;
        bannerPreview.value = URL.createObjectURL(file);
    }
};

onMounted(async () => {
    const id = sessionStorage.getItem("userId");
    const token = sessionStorage.getItem("userToken");

    if (!id || !token) return;

    try {
        const response = await fetch(`/go/users/${id}`, {
            method: "GET",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
        });

        if (response.ok) {
            const data = await response.json();
            form.value = data;

            if (data.date_naissance) {
                form.value.date_naissance = data.date_naissance.split('T')[0];
            }
            
            if (data.image_profil && data.image_profil.trim() !== "") {
                profilePreview.value = data.image_profil;
            }
            
            if (data.banniere && data.banniere.trim() !== "") {
                bannerPreview.value = data.banniere;
            }
        }
    } catch (error) {
        errors.value = ["Erreur réseau lors de la récupération du profil."];
    }
});

const formatDate = (dateString) => {
    if (!dateString) return "...";
    const date = new Date(dateString);
    return date.toLocaleDateString("fr-FR", {
        day: "numeric",
        month: "long",
        year: "numeric",
    });
};

const updateProfile = async () => {
    const userId = sessionStorage.getItem("userId");
    const token = sessionStorage.getItem("userToken");

    errors.value = [];
    successMsg.value = "";

    if (!userId || !token) return;

    try {
        const responseText = await fetch(`/go/users/${userId}`, {
            method: "PUT",
            headers: {
                Authorization: token,
                "Content-Type": "application/json",
            },
            body: JSON.stringify(form.value),
        });

        if (!responseText.ok) {
            const data = await responseText.json();
            errors.value = Array.isArray(data) ? data : [data.message || "Erreur mise à jour texte"];
            return;
        }

        if (profileFile.value || bannerFile.value) {
            const formData = new FormData();
            if (profileFile.value) formData.append("profil", profileFile.value);
            if (bannerFile.value) formData.append("banniere", bannerFile.value);

            const responseImg = await fetch(`/go/users/${userId}/images`, {
                method: "POST",
                headers: {
                    Authorization: token,
                },
                body: formData,
            });

            if (!responseImg.ok) {
                 errors.value.push("Le texte a été mis à jour, mais l'envoi des images a échoué.");
                 return;
            }
        }

        successMsg.value = "Profil mis à jour avec succès!";
        
    } catch (error) {
        errors.value = ["Le serveur est injoignable pour le moment."];
    }
};
</script>

<style scoped>
.info-view {
    padding: 10px;
}

.info-layout {
    display: grid;
    grid-template-columns: 1fr 320px;
    gap: 25px;
    margin-top: 20px;
}

.info-section {
    background: white;
    border-radius: 16px;
    padding: 24px;
    margin-bottom: 25px;
    border: 1px solid #e8ebe9;
}

.section-title {
    font-size: 1.1rem;
    font-family: "Syne", sans-serif;
    font-weight: 700;
    margin-bottom: 20px;
    color: #1a1f1c;
}

.images-row {
    display: flex;
    gap: 40px;
    align-items: flex-start;
    margin-bottom: 25px;
    flex-wrap: wrap;
}

.upload-group {
    margin-bottom: 0;
}

.upload-controls {
    display: flex;
    align-items: center;
    gap: 15px;
}

.image-label {
    display: block;
    font-weight: 600;
    color: #5a6660;
    font-size: 0.85rem;
    margin-bottom: 12px;
}

.btn-import {
    background: #f0f4f1;
    color: #2d7a4f;
    border: 1px solid #2d7a4f;
    padding: 8px 16px;
    border-radius: 8px;
    cursor: pointer;
    font-weight: 600;
    font-size: 0.85rem;
    transition: all 0.2s ease;
    white-space: nowrap;
}

.btn-import:hover {
    background: #2d7a4f;
    color: #ffffff;
}

.btn-import:active {
    transform: scale(0.95);
}

.avatar-group {
    flex-shrink: 0;
}

.avatar-container {
    width: 100px;
    height: 100px;
}

.avatar-preview {
    width: 100%;
    height: 100%;
    object-fit: cover;
    border-radius: 50%;
    border: 3px solid white;
    box-shadow: 0 4px 10px rgba(0,0,0,0.08);
}

.banner-group {
    flex-grow: 1;
    max-width: 100%;
}

.banner-preview {
    width: 100%;
    max-width: 400px;
    height: 100px;
    background-size: cover;
    background-position: center;
    border-radius: 12px;
    background-color: #f1f3f2;
}

.banner-group .upload-controls {
     justify-content: flex-start;
}

.banner-group .banner-preview {
    flex-grow: 1;
}

.input-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
}

.input-field {
    display: flex;
    flex-direction: column;
    gap: 8px;
    margin-bottom: 15px;
}

label {
    font-size: 0.85rem;
    font-weight: 600;
    color: #5a6660;
}

input {
    background-color: #fcfdfc;
    border: 1px solid #dcdfdc;
    border-radius: 10px;
    padding: 12px 15px;
    font-size: 0.95rem;
    outline: none;
    transition: border-color 0.2s;
}

input:focus {
    border-color: #2d7a4f;
}

.disabled-input {
    background-color: #f1f3f2;
    color: #888;
    cursor: not-allowed;
}

.status-item {
    display: flex;
    justify-content: space-between;
    margin-bottom: 15px;
    font-size: 0.9rem;
}

.status-label {
    color: #5a6660;
}

.badge-role {
    background-color: #eaf4ed;
    color: #2d7a4f;
    padding: 4px 10px;
    border-radius: 8px;
    font-weight: bold;
    font-size: 0.8rem;
}

.status-active {
    color: #1e7e34;
    font-weight: bold;
}

.divider {
    height: 1px;
    background: #e8ebe9;
    margin: 20px 0;
}

.registration-info p {
    font-size: 0.8rem;
    color: #5a6660;
    margin: 0;
}

.registration-info strong {
    font-size: 0.95rem;
}

.btn-outline {
    width: 100%;
    background: white;
    border: 1px solid #dcdfdc;
    padding: 10px;
    border-radius: 10px;
    cursor: pointer;
    font-weight: 600;
    transition: background 0.2s;
}

.btn-outline:hover {
    background: #f7f9f7;
}

.warning-text {
    font-size: 0.75rem;
    color: #a0ada7;
    margin-top: 10px;
    text-align: center;
}

@media (max-width: 1024px) {
    .info-layout {
        grid-template-columns: 1fr;
    }
    .info-side-col {
        order: -1;
    }
}

.siret-box {
    margin-top: 20px;
    padding-top: 20px;
    border-top: 1px dashed #dcdfdc;
}

.siret-label {
    display: flex;
    justify-content: space-between;
    align-items: center;
}

.badge-valid {
    background-color: #eaf4ed;
    color: #2d7a4f;
    padding: 3px 8px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
}

.badge-pending {
    background-color: #fff3cd;
    color: #856404;
    padding: 3px 8px;
    border-radius: 6px;
    font-size: 0.75rem;
    font-weight: bold;
}
</style>