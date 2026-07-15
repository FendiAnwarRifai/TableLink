<template>
  <div class="space-y-6 font-sans text-black">
    <!-- Header Block -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 bg-white p-6 border border-black rounded-none">
      <div>
        <h1 class="text-xl font-bold tracking-tight uppercase text-black">Items Management</h1>
        <p class="text-xs text-neutral-600">Total: {{ total }} records found.</p>
      </div>
      <button 
        @click="openCreateModal"
        class="flex items-center justify-center gap-2 bg-black text-white hover:bg-neutral-800 font-bold px-4 py-2 border border-black rounded-none transition-none"
      >
        <PlusIcon class="w-4 h-4" />
        Add Item
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
              <th class="px-6 py-3 border-r border-black">Price</th>
              <th class="px-6 py-3 border-r border-black">Ingredients</th>
              <th class="px-6 py-3 border-r border-black">Status</th>
              <th class="px-6 py-3 text-right">Actions</th>
            </tr>
          </thead>
          <tbody class="divide-y divide-black">
            <tr v-if="loading" class="text-neutral-500 text-center">
              <td colspan="5" class="py-10">LOADING...</td>
            </tr>
            <tr v-else-if="items.length === 0" class="text-neutral-500 text-center">
              <td colspan="5" class="py-10">NO DATA FOUND.</td>
            </tr>
            <tr v-else v-for="item in items" :key="item.uuid" class="hover:bg-neutral-50">
              <td class="px-6 py-4 font-bold border-r border-black text-black">{{ item.name }}</td>
              <td class="px-6 py-4 border-r border-black text-neutral-800">{{ formatPrice(item.price) }}</td>
              <td class="px-6 py-4 border-r border-black">
                <div class="flex flex-wrap gap-1.5 max-w-lg">
                  <span 
                    v-for="ing in item.ingredients" 
                    :key="ing.uuid"
                    class="inline-flex items-center px-2 py-0.5 border border-neutral-300 text-xs text-neutral-700"
                  >
                    {{ ing.name }}
                  </span>
                  <span v-if="!item.ingredients || item.ingredients.length === 0" class="text-xs text-neutral-400 italic">
                    NO INGREDIENTS
                  </span>
                </div>
              </td>
              <td class="px-6 py-4 border-r border-black">
                <span 
                  class="px-2 py-0.5 border"
                  :class="item.status === 1 ? 'border-black text-black' : 'border-neutral-300 text-neutral-400'"
                >
                  {{ item.status === 1 ? 'ACTIVE' : 'INACTIVE' }}
                </span>
              </td>
              <td class="px-6 py-4 text-right">
                <div class="flex items-center justify-end gap-3">
                  <button 
                    @click="openEditModal(item)"
                    class="px-2.5 py-1 border border-black hover:bg-neutral-100 text-xs font-bold uppercase transition-none"
                  >
                    Edit
                  </button>
                  <button 
                    @click="deleteItem(item)"
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
          <h3 class="text-sm font-bold uppercase tracking-wider text-black">{{ modal.isEdit ? 'Edit Item' : 'Add Item' }}</h3>
          <button @click="closeModal" class="text-black hover:text-neutral-600 font-bold text-xl">&times;</button>
        </div>
        
        <form @submit.prevent="saveItem" class="p-6 space-y-4 max-h-[80vh] overflow-y-auto">
          <!-- Name -->
          <div class="space-y-1">
            <label class="text-xs font-bold uppercase text-neutral-600">Name</label>
            <input 
              v-model="modal.form.name"
              type="text"
              required
              placeholder="Item name..."
              class="w-full bg-white border border-black rounded-none px-3 py-2 text-black placeholder-neutral-400 focus:outline-none"
            />
          </div>

          <!-- Price -->
          <div class="space-y-1">
            <label class="text-xs font-bold uppercase text-neutral-600">Price (IDR)</label>
            <input 
              v-model.number="modal.form.price"
              type="number"
              min="1"
              required
              placeholder="Item price..."
              class="w-full bg-white border border-black rounded-none px-3 py-2 text-black placeholder-neutral-400 focus:outline-none"
            />
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

          <!-- Ingredients Multi-Select -->
          <div class="space-y-2">
            <label class="text-xs font-bold uppercase text-neutral-600 block">Select Ingredients</label>
            <div class="border border-black bg-white rounded-none p-4 max-h-48 overflow-y-auto space-y-2">
              <div v-if="allIngredients.length === 0" class="text-neutral-400 text-xs italic">
                No active ingredients found. Create one first!
              </div>
              <div 
                v-else 
                v-for="ing in allIngredients" 
                :key="ing.uuid"
                class="flex items-center gap-3 cursor-pointer"
              >
                <input 
                  type="checkbox"
                  :id="'ing_cb_' + ing.uuid"
                  :value="ing.uuid"
                  v-model="modal.form.ingredientUuids"
                  class="w-4 h-4 rounded-none bg-white border border-black text-black focus:ring-0 focus:ring-offset-0"
                />
                <label 
                  :for="'ing_cb_' + ing.uuid"
                  class="text-xs font-bold text-neutral-700 cursor-pointer select-none uppercase"
                >
                  {{ ing.name }}
                  <span class="text-neutral-550 text-[10px] ml-1">
                    ({{ ing.cause_alergy ? 'Allergy' : 'Safe' }}, {{ ing.type === 1 ? 'Veggie' : ing.type === 2 ? 'Vegan' : 'None' }})
                  </span>
                </label>
              </div>
            </div>
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
const items = ref([]);
const allIngredients = ref([]);
const total = ref(0);
const page = ref(1);
const limit = ref(10);
const search = ref('');
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
    price: 0,
    status: 1,
    ingredientUuids: [],
  },
});

// Methods
const fetchItems = async () => {
  loading.value = true;
  error.value = null;
  try {
    const res = await axios.get(`${API_BASE}/items`, {
      params: {
        page: page.value,
        limit: limit.value,
        search: search.value || undefined,
      },
    });
    items.value = res.data.data;
    total.value = res.data.total;
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to fetch items.';
  } finally {
    loading.value = false;
  }
};

const fetchActiveIngredients = async () => {
  try {
    const res = await axios.get(`${API_BASE}/ingredients`, {
      params: {
        page: 1,
        limit: 100, // Load active ingredients
        status: 1,
      },
    });
    allIngredients.value = res.data.data;
  } catch (err) {
    console.error('Failed to load ingredients for multi-select', err);
  }
};

const onSearchChange = () => {
  page.value = 1;
  fetchItems();
};

const onLimitChange = () => {
  page.value = 1;
  fetchItems();
};

const prevPage = () => {
  if (page.value > 1) {
    page.value--;
    fetchItems();
  }
};

const nextPage = () => {
  if (page.value * limit.value < total.value) {
    page.value++;
    fetchItems();
  }
};

const formatPrice = (val) => {
  return new Intl.NumberFormat('id-ID', {
    style: 'currency',
    currency: 'IDR',
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(val);
};

// Modal Operations
const openCreateModal = async () => {
  await fetchActiveIngredients();
  modal.isEdit = false;
  modal.uuid = '';
  modal.form.name = '';
  modal.form.price = 0;
  modal.form.status = 1;
  modal.form.ingredientUuids = [];
  modal.show = true;
};

const openEditModal = async (item) => {
  await fetchActiveIngredients();
  modal.isEdit = true;
  modal.uuid = item.uuid || '';
  modal.form.name = item.name;
  modal.form.price = item.price;
  modal.form.status = item.status;
  modal.form.ingredientUuids = item.ingredients?.map(i => i.uuid) || [];
  modal.show = true;
};

const closeModal = () => {
  modal.show = false;
};

const saveItem = async () => {
  if (modal.form.price <= 0) {
    error.value = 'Price must be greater than 0.';
    return;
  }
  if (modal.form.ingredientUuids.length === 0) {
    error.value = 'At least one ingredient must be selected.';
    return;
  }

  saving.value = true;
  error.value = null;
  try {
    const payload = {
      name: modal.form.name,
      price: modal.form.price,
      status: modal.form.status,
      ingredients: modal.form.ingredientUuids,
    };

    if (modal.isEdit) {
      await axios.put(`${API_BASE}/items/${modal.uuid}`, payload);
    } else {
      await axios.post(`${API_BASE}/items`, payload);
    }
    
    closeModal();
    fetchItems();
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to save item.';
  } finally {
    saving.value = false;
  }
};

const deleteItem = async (item) => {
  if (!item.uuid) return;
  if (!confirm(`Are you sure you want to delete ${item.name}?`)) return;

  error.value = null;
  try {
    await axios.delete(`${API_BASE}/items/${item.uuid}`);
    fetchItems();
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to delete item.';
  }
};

onMounted(() => {
  fetchItems();
});
</script>
