<template>
  <div class="min-h-screen bg-white text-black flex flex-col md:flex-row font-sans">
    <!-- Sidebar -->
    <aside class="w-full md:w-64 bg-white border-b md:border-b-0 md:border-r border-black flex flex-col shrink-0">
      <!-- Navigation Links -->
      <nav class="flex-1 p-4 space-y-2">
        <button 
          @click="currentView = 'ingredients'"
          class="w-full flex items-center gap-3 px-4 py-3 rounded-none font-bold border border-black text-left transition-none"
          :class="currentView === 'ingredients' ? 'bg-black text-white border-black' : 'text-neutral-600 hover:text-black hover:bg-neutral-100'"
        >
          <CookieIcon class="w-5 h-5" />
          Ingredients
        </button>

        <button 
          @click="currentView = 'items'"
          class="w-full flex items-center gap-3 px-4 py-3 rounded-none font-bold border border-black text-left transition-none"
          :class="currentView === 'items' ? 'bg-black text-white border-black' : 'text-neutral-600 hover:text-black hover:bg-neutral-100'"
        >
          <ShoppingBagIcon class="w-5 h-5" />
          Items
        </button>
      </nav>
    </aside>

    <!-- Main Content -->
    <main class="flex-1 p-6 md:p-8 overflow-y-auto bg-white">
      <IngredientsView v-if="currentView === 'ingredients'" />
      <ItemsView v-if="currentView === 'items'" />
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import IngredientsView from './presentation/views/IngredientsView.vue';
import ItemsView from './presentation/views/ItemsView.vue';
import { Cookie as CookieIcon, ShoppingBag as ShoppingBagIcon } from '@lucide/vue';

const currentView = ref('ingredients');
</script>
