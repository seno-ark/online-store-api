package constant

import "time"

var (
	DefaultExpiration = time.Minute * 5

	CacheKeyListCategory      string = "list_category:limit_%d:offset_%d"
	CacheKeyListCategoryCount string = "list_category:count"
	CacheKeyListCategoryAll   string = "list_category:*"

	CacheKeyListProductByCategory      string = "list_product:category_%s:limit_%d:offset_%d"
	CacheKeyListProductByCategoryCount string = "list_product:category_%s:count"
	CacheKeyListProductByCategoryAll   string = "list_product:category_%s:*"

	CacheKeyListUserCart      string = "list_cart:user_%s:limit_%d:offset_%d"
	CacheKeyListUserCartCount string = "list_cart:user_%s:count"
	CacheKeyListUserCartAll   string = "list_cart:user_%s:*"

	CacheKeyListUserOrder      string = "list_order:user_%s:limit_%d:offset_%d"
	CacheKeyListUserOrderCount string = "list_order:user_%s:count"
	CacheKeyListUserOrderAll   string = "list_order:user_%s:*"

	CacheKeyListOrderItem      string = "list_order_item:order_%s:limit_%d:offset_%d"
	CacheKeyListOrderItemCount string = "list_order_item:order_%s:count"
	CacheKeyListOrderItemAll   string = "list_order_item:order_%s:*"
)
