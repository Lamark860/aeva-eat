import { createRouter, createWebHistory } from 'vue-router'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue')
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../views/Login.vue')
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../views/Register.vue')
  },
  {
    path: '/places',
    name: 'Places',
    component: () => import('../views/Places.vue')
  },
  {
    path: '/places/new',
    name: 'PlaceCreate',
    component: () => import('../views/PlaceForm.vue')
  },
  {
    path: '/places/:id',
    name: 'PlaceDetail',
    component: () => import('../views/PlaceDetail.vue')
  },
  {
    path: '/places/:id/edit',
    name: 'PlaceEdit',
    component: () => import('../views/PlaceForm.vue')
  },
  {
    path: '/profile',
    name: 'Profile',
    component: () => import('../views/Profile.vue')
  },
  {
    path: '/map',
    name: 'Map',
    component: () => import('../views/MapPage.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router
