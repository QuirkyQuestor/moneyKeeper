import { Component, OnInit, Inject } from "@angular/core";
import { Category } from "../category";
import { CategoryService } from "../category.service";
import {
  MatDialog,
  MAT_DIALOG_DATA,
  MatDialogRef,
} from "@angular/material/dialog";
import * as _ from "lodash-es";

@Component({
  selector: "app-category",
  templateUrl: "./category.component.html",
  styleUrls: ["./category.component.css"],
})
export class CategoryComponent implements OnInit {
  categories: Category[] = [];

  constructor(
    private categoryService: CategoryService,
    public dialog: MatDialog
  ) {}

  ngOnInit(): void {
    this.getCategories();
  }

  getCategories(): void {
    this.categoryService
      .getCategories()
      .subscribe((categories) => (this.categories = categories));
  }

  addCategory(
    name: string,
    parentId: string,
    description: string,
    expence: boolean
  ): void {
    name = name.trim();
    if (!name) {
      return;
    }
    this.categoryService
      .addCategory({ name, parentId, description, expence } as Category)
      .subscribe((category: Category) => {
        this.categories.push(category);
      });
  }

  deleteCategory(category: Category): void {
    console.log(category);
    this.categories = this.categories.filter((h) => h !== category);
    this.categoryService.deleteCategory(category.categoryId).subscribe();
  }

  updateCategory(category: Category): void {
    console.log(category);
    this.categoryService.updateCategory(category).subscribe();
  }

  displayedColumns: string[] = [
    "name",
    "parent",
    "description",
    "expence",
    "actions",
  ];

  getCategorytById(categoryId: string): Category {
    return this.categories.find(
      (category) => category.categoryId == categoryId
    ) as Category;
  }

  getCategorytParentName(category: Category) {
    let parentName: string = "";

    while (category.parentId) {
      let parent = this.getCategorytById(category.parentId);

      if (!parentName) {
        parentName = parent.name;
      } else {
        parentName = [parent.name, parentName].join(" :: ");
      }

      category = parent;
    }

    return parentName;
  }

  openAddCategoryDialog(): void {
    const dialogRef = this.dialog.open(CategoryAddComponent, {
      width: "250px",
      data: {},
    });

    dialogRef.afterClosed().subscribe((category) => {
      console.log("The dialog was closed");
      if (category) {
        this.addCategory(
          category.name,
          category.parentId,
          category.description,
          category.expence
        );
      }
    });
  }

  openEditCategoryDialog(category: Category): void {
    let ct = _.cloneDeep(category);

    const dialogRef = this.dialog.open(CategoryEditComponent, {
      width: "250px",
      data: category,
    });

    dialogRef.afterClosed().subscribe((category) => {
      console.log("The dialog was closed");
      if (category && !_.isEqual(category, ct)) {
        this.updateCategory(category);
      }
    });
  }

  openDeleteCategoryDialog(category: Category): void {
    const dialogRef = this.dialog.open(CategoryDeleteComponent, {
      width: "250px",
      data: category,
    });

    dialogRef.afterClosed().subscribe((category) => {
      console.log("The dialog was closed");
      if (category) {
        this.deleteCategory(category);
      }
    });
  }
}

@Component({
  selector: "category-add-dialog",
  templateUrl: "category.dialog.html",
})
export class CategoryAddComponent {
  constructor(
    public dialogRef: MatDialogRef<CategoryAddComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Category
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

@Component({
  selector: "category-edit-dialog",
  templateUrl: "category.dialog.html",
})
export class CategoryEditComponent {
  constructor(
    public dialogRef: MatDialogRef<CategoryEditComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Category
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}

@Component({
  selector: "category-delete-dialog",
  templateUrl: "category.delete.html",
})
export class CategoryDeleteComponent {
  constructor(
    public dialogRef: MatDialogRef<CategoryDeleteComponent>,
    @Inject(MAT_DIALOG_DATA) public data: Category
  ) {}

  onNoClick(): void {
    this.dialogRef.close();
  }
}
