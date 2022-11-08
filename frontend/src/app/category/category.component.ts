import { Component, OnInit } from "@angular/core";
import { Category } from "../category";
import { CategoryService } from "../category.service";

@Component({
  selector: "app-category",
  templateUrl: "./category.component.html",
  styleUrls: ["./category.component.css"],
})
export class CategoryComponent implements OnInit {
  categories: Category[] = [];
  constructor(private categoryService: CategoryService) {}

  ngOnInit(): void {
    this.getCategories();
  }

  getCategories(): void {
    this.categoryService
      .getCategories()
      .subscribe((categories) => (this.categories = categories));
  }
  addCategory(name: string, description: string): void {
    name = name.trim();
    if (!name) {
      return;
    }
    this.categoryService
      .addCategory({ name, description } as Category)
      .subscribe((category: Category) => {
        this.categories.push(category);
      });
  }

  deleteCategory(category: Category): void {
    console.log(category);
    this.categories = this.categories.filter((h) => h !== category);
    this.categoryService.deleteCategory(category.categoryId).subscribe();
  }

  displayedColumns: string[] = [
    "name",
    "parent",
    "description",
    "expence",
    "actions",
  ];
}
