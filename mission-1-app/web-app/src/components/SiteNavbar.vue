<template>
  <header class="site-navbar">
    <div class="site-navbar__inner">
      <div class="site-navbar__left">
        <RouterLink class="site-navbar__brand" to="/" aria-label="Accueil UpcycleConnect">
          <img :src="logoSrc" alt="UpcycleConnect" class="site-navbar__logo" />
        </RouterLink>

        <div v-if="isAuthenticated" class="site-navbar__account nav-menu">
          <button class="site-navbar__account-button" type="button" aria-haspopup="true">
            <span class="site-navbar__avatar">{{ userInitials }}</span>
            <span class="site-navbar__account-text">
              <strong>{{ userName }}</strong>
              <small>{{ activeUserRole }} · Score {{ userScore }} pts</small>
            </span>
          </button>
          
          <div class="nav-menu__panel nav-menu__panel--account">
            <RouterLink to="/profil/informations" class="nav-menu__item">
              <span>Compte</span>
              <small>Informations personnelles</small>
            </RouterLink>
            <RouterLink to="/profil" class="nav-menu__item">
              <span>Activite</span>
              <small>Score, planning et resume</small>
            </RouterLink>
            <RouterLink to="/profil/factures" class="nav-menu__item">
              <span>Factures</span>
              <small>Historique d'achats et reçus</small>
            </RouterLink>
            
            <RouterLink v-if="activeUserRole === 'Prestataire'" to="/profil/abonnement" class="nav-menu__item">
              <span>Abonnement</span>
              <small>Gestion du forfait pro et alertes</small>
            </RouterLink>

            <button class="nav-menu__item nav-menu__item--danger" type="button" @click="handleLogout">
              <span>Se deconnecter</span>
              <small>Fermer la session locale</small>
            </button>
          </div>
        </div>
      </div>

      <nav class="site-navbar__links" aria-label="Navigation principale">
        <div
          v-for="group in dynamicNavGroups"
          :key="group.label"
          class="nav-menu"
          :class="{ 'is-active': isGroupActive(group) }"
        >
          <RouterLink
            v-if="group.to"
            :to="group.to"
            class="site-navbar__link"
            :class="{ 'is-active': isGroupActive(group) }"
          >
            {{ group.label }}
          </RouterLink>
          <button
            v-else
            class="site-navbar__link site-navbar__link--button"
            :class="{ 'is-active': isGroupActive(group) }"
            type="button"
            aria-haspopup="true"
          >
            {{ group.label }}
            <span class="site-navbar__chevron">⌄</span>
          </button>

          <div v-if="group.children?.length" class="nav-menu__panel">
            <RouterLink
              v-for="item in group.children"
              :key="item.label"
              :to="item.to"
              class="nav-menu__item"
              :class="{ 'is-active': isActive(item) }"
            >
              <span>{{ item.label }}</span>
              <small>{{ item.description }}</small>
            </RouterLink>
          </div>
        </div>
      </nav>

      <div class="site-navbar__actions" v-if="isAuthenticated">
        <RouterLink class="site-navbar__button site-navbar__button--primary" to="/panier">
        Panier
        </RouterLink>
        <RouterLink v-if="activeUserRole !== 'Prestataire'" class="site-navbar__button site-navbar__button--primary" to="/profil/annonces">
          + Deposer
        </RouterLink>
      </div>

      <div class="site-navbar__actions" v-else>
        <RouterLink class="site-navbar__button site-navbar__button--ghost" to="/connexion">
          Connexion
        </RouterLink>
        <RouterLink class="site-navbar__button site-navbar__button--primary" to="/inscription">
          S'inscrire
        </RouterLink>
      </div>
      <div class="nav-menu site-navbar__langue">
        <button class="site-navbar__link site-navbar__link--button" type="button" aria-haspopup="true">
            {{ currentLangCode.toUpperCase() }}
            <span class="site-navbar__chevron">⌄</span>
        </button>
        
        <div class="nav-menu__panel nav-menu__panel--lang">
            <button 
            v-for="langue in langues" 
            :key="langue.id" 
            class="nav-menu__item nav-menu__item--lang"
            @click="changerLangue(langue)"
            type="button"
            >
            <span>{{ langue.code.toUpperCase() }}</span>
            </button>
        </div>
        </div>
    </div>
  </header>
</template>

<script setup>
import { computed, ref, onMounted } from "vue";
import { RouterLink, useRoute, useRouter } from "vue-router";
import logoSrc from "./logo_texte.png";

const route = useRoute();
const router = useRouter();

const API_URL = "http://localhost:8081";

const langues = ref([]);
const currentLangCode = ref(localStorage.getItem("langCode") || "fr");
const t = ref({});

const fetchLangues = async () => {
  try {
    const res = await fetch(`${API_URL}/langues`);
    if (res.ok) langues.value = await res.json();
  } catch (e) { console.error(e); }
};

const fetchTraductions = async (code) => {
  try {
    const res = await fetch(`${API_URL}/traductions/${code}`);
    if (res.ok) t.value = await res.json();
  } catch (e) { console.error(e); }
};

const changerLangue = async (langue) => {
  if (!langue) return;
  
  if (document.activeElement) {
    document.activeElement.blur();
  }

  currentLangCode.value = langue.code;
  localStorage.setItem("langCode", langue.code);
  localStorage.setItem("langId", langue.id);
  await fetchTraductions(langue.code);

  window.dispatchEvent(new Event("lang-changed"));

  if (props.isAuthenticated) {
    const userId = sessionStorage.getItem("userId");
    if (userId) {
      fetch(`${API_URL}/users/${userId}/langue`, {
        method: "PUT",
        headers: { 
          "Content-Type": "application/json",
          "Authorization": `Bearer ${sessionStorage.getItem("userToken")}`
        },
        body: JSON.stringify({ id_langue: langue.id })
      }).catch(e => console.error(e));
    }
  }
};

onMounted(() => {
  fetchLangues();
  fetchTraductions(currentLangCode.value);
});

const props = defineProps({
  variant: {
    type: String,
    default: "public"
  },
  userName: {
    type: String,
    default: "Marie Lambert"
  },
  userRole: {
    type: String,
    default: "Particulier"
  },
  userScore: {
    type: [String, Number],
    default: () => sessionStorage.getItem("userScore") || 0
  },
  isAuthenticated: {
    type: Boolean,
    default: false
  }
});

const activeUserRole = computed(() => {
  return sessionStorage.getItem("userRole") || props.userRole;
});

const dashboardChildrenParticulier = [
  {
    label: "Vue d'ensemble",
    to: "/profil",
    description: "Score, planning et activite"
  },
  {
    label: "Mes annonces",
    to: "/profil/annonces",
    description: "Objets publies et brouillons"
  },
  {
    label: "Mes depots conteneurs",
    to: "/profil/depots",
    description: "Suivi des depots et collectes"
  },
  {
    label: "Informations",
    to: "/profil/informations",
    description: "Profil, adresse et preferences"
  }
];

const dashboardChildrenPrestataire = [
  {
    label: "Vue d'ensemble",
    to: "/profil",
    description: "Récupération objets en conteneurs"
  },
  {
    label: "Récupérer des objets",
    to: "/profil/recuperations",
    description: "Objets et matériaux en attente de récupération"
  },
  {
    label: "Mes annonces likées",
    to: "/profil/favoris",
    description: "Objets et matériaux sauvegardés"
  },
  {
    label: "Mes projets",
    to: "/profil/projets",
    description: "Créations et vitrine d'upcycling"
  },
  {
    label: "Informations",
    to: "/profil/informations",
    description: "Profil pro et informations"
  }
];

const serviceChildren = [
  {
    label: "Upcycling Score",
    to: "/profil",
    description: "Impact et points utilisateur"
  },
  {
    label: "Formations & Ateliers",
    to: "/formations",
    description: "Sessions et apprentissage"
  },
  {
    label: "Mon planning",
    to: "/profil/planning",
    description: "Evenements et rendez-vous"
  },
  {
    label: "Espace Conseils",
    to: "/conseils",
    description: "Guides et astuces"
  }
];

const communityChildren = [
  {
    label: "Catalogue offres",
    to: "/catalogue",
    description: "Toutes les annonces publiques"
  },
  {
    label: "Forums",
    to: "/forums",
    description: "Discussions, entraide et projets"
  },
  {
    label: "Projets Upcycling",
    to: "/projets",
    description: "Rejoignez des initiatives de création"
  },
  {
    label: "Evenements",
    to: "/evenements",
    description: "Participez à des événements"
  }
];

const objetsChildren = [
  {
    label: "Depot Objet",
    to: "/deposer",
    description: "Déposez vos objets"
  },
  {
    label: "Récupérer Objet",
    to: "/claim",
    description: "Récupérez vos objets"
  },
];

const dynamicNavGroups = computed(() => {
  const groups = [
    {
      label: "Tableau de bord",
      children: activeUserRole.value === "Prestataire" ? dashboardChildrenPrestataire : dashboardChildrenParticulier
    },
    {
      label: "Services",
      children: serviceChildren
    },
    {
      label: "Communaute",
      children: communityChildren
    }
  ];

  if (activeUserRole.value === "Salarie") {
    groups.push({
      label: "Réclamer ou Retirer un objet", 
      children: objetsChildren
    });
  }

  return groups;
});

const userInitials = computed(() =>
  props.userName
    .split(" ")
    .filter(Boolean)
    .slice(0, 2)
    .map((part) => part[0]?.toUpperCase() ?? "")
    .join("")
);

function itemPath(item) {
  return typeof item.to === "string" ? item.to : item.to?.path;
}

function isActive(item) {
  const path = itemPath(item);
  if (!path) return false;
  if (path === "/" && item.to?.hash) {
    return route.path === "/" && route.hash === item.to.hash;
  }
  if (path === "/profil") {
    return route.path === "/profil";
  }
  return route.path === path || route.path.startsWith(`${path}/`);
}

function isGroupActive(group) {
  if (group.children?.length) {
    return group.children.some((item) => isActive(item));
  }
  return isActive(group);
}

function handleLogout() {
  sessionStorage.removeItem("userToken");
  sessionStorage.removeItem("userId");
  sessionStorage.removeItem("userRole"); 
  router.push("/connexion");
}
</script>

<style>
.site-navbar__right {
  display: flex;
  align-items: center;
  gap: 15px;
}

.site-navbar__langue {
  position: relative;
  margin-left: 10px;
}

.nav-menu__panel--lang {
  min-width: 60px; 
  right: 0 !important; 
  left: auto !important; 
  transform: none !important; 
  padding: 5px;
}

.nav-menu__panel--lang {
  width: 60px !important;      
  min-width: 0 !important;    
  right: 0 !important;
  left: auto !important;
  transform: none !important;
  padding: 5px;
}

.nav-menu__item--lang:hover {
  background: #f0f7f3;
  border-radius: 6px;
}
</style>