<template>
  <div class="space-y-6 font-sans text-black">
    <!-- Header Block -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 bg-white p-6 border border-black rounded-none">
      <div>
        <h1 class="text-xl font-bold tracking-tight uppercase text-black">Ingredients Management</h1>
        <p class="text-xs text-neutral-600">Total: {{ total }} records found.</p>
      </div>
      <button 
        @click="openCreateModal"
        class="flex items-center justify-center gap-2 bg-black text-white hover:bg-neutral-800 font-bold px-4 py-2 border border-black rounded-none transition-none"
      >
        <PlusIcon class="w-4 h-4" />
        Add Ingredient
      </button>
    </div>

    <!-- Error Alert -->
    <div v-if="error" class="bg-white border border-red-650 text-red-600 p-4 rounded-none flex items-center justify-between">
      <span>ERROR: {{ error }}</span>
      <button @click="error = null" class="text-red-600 hover:text-red-800 font-bold">&times;</button>
    </div>

    <!-- Table Section -->
    <div class="bg-white border border-black rounded-none overflow-hidden">
      <div class="overflow-x-auto">
        <table class="w-full text-left border-collapse">
          <thead>
            <tr class="bg-neutral-100 border-b border-black text-xs font-bold uppercase text-black">
              <th class="px-6 py-3 border-r border-black">Name</th>
              <th class="px-6 py-3 border-r border-black">Allergy Cause</th>
              <th class="px-6 py-3 border-r border-black">Type</th>
              <th class="px-6 py-3 border-r border-black">Status</th>
              <th class="px-6 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-black">
            <tr v-if="loading" class="text-neutral-500 text-center">
              <td colspan="5" class="py-10">LOADING...</td>
            </tr>
            <tr v-else-if="ingredients.length === 0" class="text-neutral-500 text-center">
              <td colspan="5" class="py-10">NO DATA FOUND.</td>
            </tr>
            <tr v-else v-for="ing in ingredients" :key="ing.uuid" class="hover:bg-neutral-50">
              <td class="px-6 py-4 font-bold border-r border-black text-black">{{ ing.name }}</td>
              <td class="px-6 py-4 border-r border-black">
                <span 
                  class="px-2 py-0.5 border"
                  :class="ing.cause_alergy ? 'border-black bg-black text-white' : 'border-neutral-300 text-neutral-500'"
                >
                  {{ ing.cause_alergy ? 'ALLERGY' : 'SAFE' }}
                </span>
              </td>
              <td class="px-6 py-4 border-r border-black text-neutral-700">
                {{ formatType(ing.type) }}
              </td>
              <td class="px-6 py-4 border-r border-black">
                <span 
                  class="px-2 py-0.5 border"
                  :class="ing.status === 1 ? 'border-black text-black' : 'border-neutral-300 text-neutral-400'"
                >
                  {{ ing.status === 1 ? 'ACTIVE' : 'INACTIVE' }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-3">
                  <button 
                    @click="openEditModal(ing)"
                    class="px-2.5 py-1 border border-black hover:bg-neutral-100 text-xs font-bold uppercase transition-none"
                  >
                    Edit
                  </button>
                  <button 
                    @click="deleteIngredient(ing)"
                    class="px-2.5 py-1 border border-black hover:bg-neutral-100 text-xs font-bold uppercase text-neutral-600 hover:text-black transition-none"
                  >
                    Delete
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination Footer -->
      <div class="bg-white px-6 py-4 border-t border-black flex flex-col sm:flex-row items-center justify-between gap-4">
        <div class="flex items-center gap-2">
          <span class="text-xs text-neutral-600 uppercase font-bold">Show</span>
          <select 
            v-model="limit"
            @change="onLimitChange"
            class="bg-white border border-black rounded-none px-2 py-1 text-xs text-black focus:outline-none"
          >
            <option :value="10">10</option>
            <option :value="20">20</option>
            <option :value="50">50</option>
          </select>
          <span class="text-xs text-neutral-600 uppercase font-bold">per page</span>
        </div>

        <div class="flex items-center gap-4">
          <span class="text-xs text-neutral-600 font-bold uppercase">
            Showing {{ total === 0 ? 0 : (page - 1) * limit + 1 }} to {{ Math.min(page * limit, total) }} of {{ total }} entries
          </span>
          <div class="flex gap-2">
            <button 
              @click="prevPage"
              :disabled="page === 1"
              class="px-3 py-1 bg-white hover:bg-neutral-100 disabled:opacity-30 border border-black disabled:hover:bg-white text-xs font-bold uppercase transition-none"
            >
              Prev
            </button>
            <button 
              @click="nextPage"
              :disabled="page * limit >= total"
              class="px-3 py-1 bg-white hover:bg-neutral-100 disabled:opacity-30 border border-black disabled:hover:bg-white text-xs font-bold uppercase transition-none"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="modal.show" class="fixed inset-0 z-50 flex items-center justify-center p-4">
      <!-- Backdrop -->
      <div class="absolute inset-0 bg-black/55" @click="closeModal"></div>
      
      <!-- Modal Box -->
      <div class="relative bg-white border-2 border-black w-full max-w-md rounded-none z-10 font-mono text-black">
        <div class="px-6 py-4 border-b border-black flex items-center justify-between">
          <h3 class="text-sm font-bold uppercase tracking-wider text-black">{{ modal.isEdit ? 'Edit Ingredient' : 'Add Ingredient' }}</h3>
          <button @click="closeModal" class="text-black hover:text-neutral-600 font-bold text-xl">&times;</button>
        </div>
        
        <form @submit.prevent="saveIngredient" class="p-6 space-y-4">
          <!-- Name -->
          <div class="space-y-1">
            <label class="text-xs font-bold uppercase text-neutral-600">Name</label>
            <input 
              v-model="modal.form.name"
              type="text"
              required
              placeholder="Ingredient name..."
              class="w-full bg-white border border-black rounded-none px-3 py-2 text-black placeholder-neutral-400 focus:outline-none"
            />
          </div>

          <!-- Allergy Cause -->
          <div class="space-y-1">
            <label class="text-xs font-bold uppercase text-neutral-600">Allergy Cause</label>
            <select 
              v-model="modal.form.causeAlergy"
              class="w-full bg-white border border-black rounded-none px-3 py-2 text-black focus:outline-none"
            >
              <option :value="false">No (Safe)</option>
              <option :value="true">Yes (Allergy)</option>
            </select>
          </div>

          <!-- Diet Type -->
          <div class="space-y-1">
            <label class="text-xs font-bold uppercase text-neutral-600">Diet Type</label>
            <select 
              v-model="modal.form.type"
              class="w-full bg-white border border-black rounded-none px-3 py-2 text-black focus:outline-none"
            >
              <option :value="0">None</option>
              <option :value="1">Veggie</option>
              <option :value="2">Vegan</option>
            </select>
          </div>

          <!-- Status -->
          <div class="space-y-1">
            <label class="text-xs font-bold uppercase text-neutral-600">Status</label>
            <select 
              v-model="modal.form.status"
              class="w-full bg-white border border-black rounded-none px-3 py-2 text-black focus:outline-none"
            >
              <option :value="1">Active</option>
              <option :value="0">Inactive</option>
            </select>
          </div>

          <!-- Footer Actions -->
          <div class="flex items-center justify-end gap-3 pt-4 border-t border-black">
            <button 
              type="button"
              @click="closeModal"
              class="px-4 py-2 border border-black hover:bg-neutral-100 bg-white text-xs font-bold uppercase transition-none text-black"
            >
              Cancel
            </button>
            <button 
              type="submit"
              :disabled="saving"
              class="px-4 py-2 bg-black text-white hover:bg-neutral-800 border border-black disabled:opacity-50 text-xs font-bold uppercase transition-none"
            >
              {{ saving ? 'Saving...' : 'Save' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, reactive } from 'vue';
import axios from 'axios';
import { Plus as PlusIcon } from '@lucide/vue';

const API_BASE = 'http://localhost:3000/api';

// State
const ingredients = ref([]);
const total = ref(0);
const page = ref(1);
const limit = ref(10);
const search = ref('');
const allergyFilter = ref('');
const typeFilter = ref('');
const statusFilter = ref('');
const loading = ref(false);
const saving = ref(false);
const error = ref(null);

// Modal State
const modal = reactive({
  show: false,
  isEdit: false,
  uuid: '',
  form: {
    name: '',
    causeAlergy: false,
    type: 0,
    status: 1,
  },
});

// Methods
const fetchIngredients = async () => {
  loading.value = true;
  error.value = null;
  try {
    const res = await axios.get(`${API_BASE}/ingredients`, {
      params: {
        page: page.value,
        limit: limit.value,
        search: search.value || undefined,
        cause_alergy: allergyFilter.value !== '' ? allergyFilter.value : undefined,
        type: typeFilter.value !== '' ? typeFilter.value : undefined,
        status: statusFilter.value !== '' ? statusFilter.value : undefined,
      },
    });
    ingredients.value = res.data.data;
    total.value = res.data.total;
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to fetch ingredients.';
  } finally {
    loading.value = false;
  }
};

const onSearchChange = () => {
  page.value = 1;
  fetchIngredients();
};

const onFilterChange = () => {
  page.value = 1;
  fetchIngredients();
};

const onLimitChange = () => {
  page.value = 1;
  fetchIngredients();
};

const prevPage = () => {
  if (page.value > 1) {
    page.value--;
    fetchIngredients();
  }
};

const nextPage = () => {
  if (page.value * limit.value < total.value) {
    page.value++;
    fetchIngredients();
  }
};

const formatType = (val) => {
  if (val === 1) return 'VEGGIE';
  if (val === 2) return 'VEGAN';
  return 'NONE';
};

// Modal Operations
const openCreateModal = () => {
  modal.isEdit = false;
  modal.uuid = '';
  modal.form.name = '';
  modal.form.causeAlergy = false;
  modal.form.type = 0;
  modal.form.status = 1;
  modal.show = true;
};

const openEditModal = (ing) => {
  modal.isEdit = true;
  modal.uuid = ing.uuid || '';
  modal.form.name = ing.name;
  modal.form.causeAlergy = ing.cause_alergy;
  modal.form.type = ing.type;
  modal.form.status = ing.status;
  modal.show = true;
};

const closeModal = () => {
  modal.show = false;
};

const saveIngredient = async () => {
  saving.value = true;
  error.value = null;
  try {
    const payload = {
      name: modal.form.name,
      cause_alergy: modal.form.causeAlergy,
      type: modal.form.type,
      status: modal.form.status,
    };

    if (modal.isEdit) {
      await axios.put(`${API_BASE}/ingredients/${modal.uuid}`, payload);
    } else {
      await axios.post(`${API_BASE}/ingredients`, payload);
    }
    
    closeModal();
    fetchIngredients();
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to save ingredient.';
  } finally {
    saving.value = false;
  }
};

const deleteIngredient = async (ing) => {
  if (!ing.uuid) return;
  if (!confirm(`Are you sure you want to delete ${ing.name}?`)) return;

  error.value = null;
  try {
    await axios.delete(`${API_BASE}/ingredients/${ing.uuid}`);
    fetchIngredients();
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to delete ingredient.';
  }
};

onMounted(() => {
  fetchIngredients();
});
</script>
