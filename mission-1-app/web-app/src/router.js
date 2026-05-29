import { createRouter, createWebHistory } from "vue-router";

import Register from "./pages/Register.vue";
import Login from "./pages/Login.vue";
import Profil from "./pages/Profil.vue";

import Home from "./pages/Home.vue";
import Annonces from "./pages/Annonces.vue";
import Forums from "./pages/Forums.vue";

import Formations from "./pages/Formations.vue";
import Evenements from "./pages/Evenements.vue";
import Projets from "./pages/Projets.vue";
import Conseils from "./pages/Conseils.vue";
import Panier from "./pages/PanierDetail.vue";
import ClaimObjet from "./pages/ClaimObjet.vue";

const routes = [
  { path: "/", component: Home },
  { path: "/inscription", component: Register },
  { path: "/connexion", component: Login },
  { path: "/catalogue", component: Annonces },
  { path: "/forums", component: Forums },
  { path: "/formations", component: Formations },
  { path: "/evenements", component: Evenements },
  { path: "/projets", component: Projets },
  { path: "/conseils", component: Conseils },
  { path: "/panier", component: Panier },
  {
    path: '/annonce/:id',
    name: 'annonce-detail',
    component: () => import('./pages/AnnonceDetail.vue')
  },
  {
    path: '/user/:id',
    name: 'utilisateur-detail',
    component: () => import('./pages/UtilisateurDetail.vue')
  },
  {
    path: '/formations/:id',
    name: 'formation-detail',
    component: () => import('./pages/FormationDetail.vue')
  },
  {
    path: '/projets/:id',
    name: 'projet-detail',
    component: () => import('./pages/ProjetDetail.vue')
  },
  {
    path: '/conseils/:id',
    name: 'conseil-detail',
    component: () => import('./pages/ConseilDetail.vue')
  },
  {
    path: "/deposer",
    name: "deposer-objet",
    component: () => import('./pages/DropObjet.vue'),
    meta: { requiresRole: "Salarie" }
  },
  {
    path: "/claim",
    name: "recuperer-objet",
    component: () => import('./pages/ClaimObjet.vue'),
    meta: { requiresRole: "Salarie" }
  },
  {
    path: "/profil",
    component: Profil,
    meta: { requiresAuth: true },
    children: [
      {
        path: "",
        name: "dashboard",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/DashboardHome.vue');
          if (role === "Salarie") return import('./pages/profilpagesSalarie/DashboardHome.vue');
          return import('./pages/profilpagesParticulier/DashboardHome.vue');
        }
      },
      {
        path: "planning",
        name: "mon-planning",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/Planning.vue');
          if (role === "Salarie") return import('./pages/profilpagesSalarie/Planning.vue');
          return import('./pages/profilpagesParticulier/Planning.vue');
        }
      },
      {
        path: "informations",
        name: "mes-informations",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/Infos.vue');
          if (role === "Salarie") return import('./pages/profilpagesSalarie/Infos.vue');
          return import('./pages/profilpagesParticulier/Infos.vue');
        }
      },
      {
        path: "factures",
        name: "mes-factures",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/Factures.vue');
          return import('./pages/profilpagesParticulier/Factures.vue');
        }
      },
      {
        path: "forums",
        name: "mes-forums",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Salarie") return import('./pages/profilpagesSalarie/Forums.vue');
          return import('./pages/profilpagesParticulier/Forums.vue');
        }
      },
      {
        path: "password",
        name: "modification-password",
        component: () => {
          const role = sessionStorage.getItem("userRole");
          if (role === "Prestataire") return import('./pages/profilpagesPrestataire/ChangePassword.vue');
          if (role === "Salarie") return import('./pages/profilpagesSalarie/ChangePassword.vue');
          return import('./pages/profilpagesParticulier/ChangePassword.vue');
        }
      },
      
      {
        path: "annonces",
        name: "mes-annonces",
        component: () => import('./pages/profilpagesParticulier/Annonces.vue'),
        meta: { requiresRole: "Particulier" }
      },
      {
        path: "depots",
        name: "mes-depots",
        component: () => import('./pages/profilpagesParticulier/Depots.vue'),
        meta: { requiresRole: "Particulier" }
      },
      {
        path: "modifyAnnonce/:id",
        name: "modification-annonce",
        component: () => import('./pages/profilpagesParticulier/modifyAnnonce.vue'),
        meta: { requiresRole: "Particulier" }
      },
      {
        path: "createAnnonce",
        name: "create-annonce",
        component: () => import('./pages/profilpagesParticulier/CreateAnnonce.vue'),
        meta: { requiresRole: "Particulier" }
      },
      {
        path: "reserveCasier/:id",
        name: "reserve-casier",
        component: () => import('./pages/profilpagesParticulier/ReserverCasier.vue'),
        meta: { requiresRole: "Particulier" }
      },

      // --- SPÉCIFIQUES AUX PRESTATAIRES ---
      {
        path: "favoris",
        name: "annonce-like",
        component: () => import("./pages/profilpagesPrestataire/AnnoncesLike.vue"),
        meta: { requiresRole: "Prestataire" }
      },
      {
        path: "recuperations",
        name: "objet-recuperation",
        component: () => import("./pages/profilpagesPrestataire/ObjetRecuperation.vue"),
        meta: { requiresRole: "Prestataire" }
      },
      {
        path: "projets",
        name: "mes-projet",
        component: () => import("./pages/profilpagesPrestataire/Projets.vue"),
        meta: { requiresRole: "Prestataire" }
      },
      {
        path: "modifyProjet/:id",
        name: "modify-projet",
        component: () => import("./pages/profilpagesPrestataire/ModifyProjet.vue"),
        meta: { requiresRole: "Prestataire" }
      },
      {
        path: "createProjet",
        name: "create-projet",
        component: () => import("./pages/profilpagesPrestataire/CreateProjet.vue"),
        meta: { requiresRole: "Prestataire" }
      },
      {
        path: "abonnement",
        name: "mon-abonnement",
        component: () => import("./pages/profilpagesPrestataire/Abonnement.vue"),
        meta: { requiresRole: "Prestataire" }
      }
    ],
  },
];

const router = createRouter({
  history: createWebHistory(),
  routes,
  scrollBehavior(to, from, savedPosition) {
    if (to.hash) {
      return { el: to.hash, behavior: "smooth" };
    }
    return { top: 0 };
  },
});

router.beforeEach(async (to, from, next) => {
  const token = sessionStorage.getItem("userToken");
  const id = sessionStorage.getItem("userId");
  const userRole = sessionStorage.getItem("userRole");

  // 1. Vérification d'authentification globale
  if (to.matched.some((record) => record.meta.requiresAuth)) {
    if (!id || !token) {
      return next("/connexion");
    }

    try {
      const response = await fetch(
        `http://localhost:8081/check-session?id=${id}`,
        {
          headers: { Authorization: token },
        }
      );
      const data = await response.json();

      if (!data.isValid) {
        sessionStorage.clear();
        return next("/connexion");
      }
    } catch (error) {
      return next("/");
    }
  }

  // 2. Vérification stricte des autorisations de rôle
  const requiredRole = to.matched.find(record => record.meta.requiresRole)?.meta.requiresRole;
  
  if (requiredRole && userRole !== requiredRole) {
    console.warn(`Accès refusé ! Rôle attendu : ${requiredRole}, Rôle actuel : ${userRole}`);
    
    if (userRole === "Prestataire" || userRole === "Salarie" || userRole === "Particulier") {
      return next("/profil"); 
    } else {
      return next("/"); 
    }
  }

  next();
});

export default router;